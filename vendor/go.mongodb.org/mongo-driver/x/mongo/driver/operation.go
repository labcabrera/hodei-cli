package driver

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/golang/snappy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	wiremessagex "go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
	"go.mongodb.org/mongo-driver/x/mongo/driverlegacy/session"
	"go.mongodb.org/mongo-driver/x/network/description"
	"go.mongodb.org/mongo-driver/x/network/wiremessage"
)

var dollarCmd = [...]byte{'.', '$', 'c', 'm', 'd'}

var (
	// ErrNoDocCommandResponse occurs when the server indicated a response existed, but none was found.
	ErrNoDocCommandResponse = errors.New("command returned no documents")
	// ErrMultiDocCommandResponse occurs when the server sent multiple documents in response to a command.
	ErrMultiDocCommandResponse = errors.New("command returned multiple documents")
)

// InvalidOperationError is returned from Validate and indicates that a required field is missing
// from an instance of Operation.
type InvalidOperationError struct{ MissingField string }

func (err InvalidOperationError) Error() string {
	return "the " + err.MissingField + " field must be set on Operation"
}

// Operation is used to execute an operation. It contains all of the common code required to
// select a server, transform an operation into a command, write the command to a connection from
// the selected server, read a response from that connection, process the response, and potentially
// retry.
//
// The required fields are Database, CommandFn, and Deployment. All other fields are optional.
//
// While an Operation can be constructed manually, drivergen should be used to generate an
// implementation of an operation instead. This will ensure that there are helpers for constructing
// the operation and that this type isn't configured incorrectly.
type Operation struct {
	// CommandFn is used to create the command that will be wrapped in a wire message and sent to
	// the server. This function should only add the elements of the command and not start or end
	// the enclosing BSON document. Per the command API, the first element must be the name of the
	// command to run. This field is required.
	CommandFn func(dst []byte, desc description.SelectedServer) ([]byte, error)

	// Database is the database that the command will be run against. This field is required.
	Database string

	// Deployment is the MongoDB Deployment to use. While most of the time this will be multiple
	// servers, commands that need to run against a single, preselected server can use the
	// SingleServerDeployment type. Commands that need to run on a preselected connection can use
	// the SingleConnectionDeployment type.
	Deployment Deployment

	// ProcessResponseFn is called after a response to the command is returned. The server is
	// provided for types like Cursor that are required to run subsequent commands using the same
	// server.
	ProcessResponseFn func(response bsoncore.Document, srvr Server) error

	// Selector is the server selector that's used during both initial server selection and
	// subsequent selection for retries. Depending on the Deployment implementation, the
	// SelectServer method may not actually be called.
	Selector description.ServerSelector

	// ReadPreference is the read preference that will be attached to the command. If this field is
	// not specified a default read preference of primary will be used.
	ReadPreference *readpref.ReadPref

	// ReadConcern is the read concern used when running read commands. This field should not be set
	// for write operations. If this field is set, it will be encoded onto the commands sent to the
	// server.
	ReadConcern *readconcern.ReadConcern

	// WriteConcern is the write concern used when running write commands. This field should not be
	// set for read operations. If this field is set, it will be encoded onto the commands sent to
	// the server.
	WriteConcern *writeconcern.WriteConcern

	// Client is the session used with this operation. This can be either an implicit or explicit
	// session. If the server selected does not support sessions and Client is specified the
	// behavior depends on the session type. If the session is implicit, the session fields will not
	// be encoded onto the command. If the session is explicit, an error will be returned. The
	// caller is responsible for ensuring that this field is nil if the Deployment does not support
	// sessions.
	Client *session.Client

	// Clock is a cluster clock, different from the one contained within a session.Client. This
	// allows updating cluster times for a global cluster clock while allowing individual session's
	// cluster clocks to be only updated as far as the last command that's been run.
	Clock *session.ClusterClock

	// RetryMode specifies how to retry. There are three modes that enable retry: RetryOnce,
	// RetryOncePerCommand, and RetryContext. For more information about what these modes do, please
	// refer to their definitions. Both RetryMode and RetryType must be set for retryability to be
	// enabled.
	RetryMode *RetryMode

	// RetryType specifies the kinds of operations that can be retried. There is only one mode that
	// enables retry: RetryWrites. For more information about what this mode does, please refer to
	// it's definition. Both RetryType and RetryMode must be set for retryability to be enabled.
	RetryType RetryType

	// Batches contains the documents that are split when executing a write command that potentially
	// has more documents than can fit in a single command. This should only be specified for
	// commands that are batch compatible. For more information, please refer to the definition of
	// Batches.
	Batches *Batches

	// Legacy sets the legacy type for this operation. There are only 3 types that require legacy
	// support: find, getMore, and killCursors. For more information about LegacyOperationKind,
	// please refer to it's definition.
	Legacy LegacyOperationKind
}

// selectServer handles performing server selection for an operation.
func (op Operation) selectServer(ctx context.Context) (Server, error) {
	if err := op.Validate(); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	selector := op.Selector
	if selector == nil {
		rp := op.ReadPreference
		if rp == nil {
			rp = readpref.Primary()
		}
		selector = description.CompositeSelector([]description.ServerSelector{
			description.ReadPrefSelector(rp),
			description.LatencySelector(15 * time.Millisecond),
		})
	}

	return op.Deployment.SelectServer(ctx, selector)
}

// Validate validates this operation, ensuring the fields are set properly.
func (op Operation) Validate() error {
	if op.CommandFn == nil {
		return InvalidOperationError{MissingField: "CommandFn"}
	}
	if op.Deployment == nil {
		return InvalidOperationError{MissingField: "Deployment"}
	}
	if op.Database == "" {
		return InvalidOperationError{MissingField: "Database"}
	}
	return nil
}

// Execute runs this operation. The scratch parameter will be used and overwritten (potentially many
// times), this should mainly be used to enable pooling of byte slices.
func (op Operation) Execute(ctx context.Context, scratch []byte) error {
	err := op.Validate()
	if err != nil {
		return err
	}

	srvr, err := op.selectServer(ctx)
	if err != nil {
		return err
	}

	conn, err := srvr.Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	desc := description.SelectedServer{Server: conn.Description(), Kind: op.Deployment.Kind()}

	// TODO(GODRIVER-617): We should check the wire version here. If we're doing a find, getMore, or
	// killCursors and the wire version is less than 4 we need to call out to legacy code here.
	if desc.WireVersion == nil || desc.WireVersion.Max < 4 {
		switch op.Legacy {
		case LegacyFind:
			// TODO(GODRIVER-984): Implement LegacyFind.
			return errors.New("legacy find is not yet supported")
		case LegacyGetMore:
			// TODO(GODRIVER-984): Implement LegacyGetMore.
			return errors.New("legacy getMore is not yet supported")
		case LegacyKillCursors:
			// TODO(GODRIVER-984): Implement LegacyKillCursors.
			return errors.New("legacy killCursors is not yet supported")
		}
	}

	var res bsoncore.Document
	var operationErr WriteCommandError
	var original error
	var retries int
	// TODO(GODRIVER-617): Add support for retryable reads.
	retryable := op.retryable(desc.Server)
	if retryable == RetryWrite && op.Client != nil && op.RetryMode != nil {
		if *op.RetryMode > RetryNone {
			op.Client.RetryWrite = true
			op.Client.IncrementTxnNumber()
		}

		switch *op.RetryMode {
		case RetryOnce, RetryOncePerCommand:
			retries = 1
		case RetryContext:
			retries = -1
		}
	}
	batching := op.Batches.Valid()
	for {
		if batching {
			err = op.Batches.AdvanceBatch(int(desc.MaxBatchCount), int(desc.MaxDocumentSize))
			if err != nil {
				// TODO(GODRIVER-982): Should we also be returning operationErr?
				return err
			}
		}

		// convert to wire message
		if len(scratch) > 0 {
			scratch = scratch[:0]
		}
		wm, err := op.createWireMessage(scratch, desc)
		if err != nil {
			return err
		}

		// compress wiremessage if allowed
		if compressor, ok := conn.(Compressor); ok && op.canCompress("") {
			wm, err = compressor.CompressWireMessage(wm, nil)
			if err != nil {
				return err
			}
		}

		// roundtrip
		wm, err = op.roundTrip(ctx, conn, wm)
		if ep, ok := srvr.(ErrorProcessor); ok {
			ep.ProcessError(err)
		}
		if err != nil {
			return err
		}

		// decompress wiremessage
		wm, err = op.decompressWireMessage(wm)
		if err != nil {
			return err
		}

		// decode
		res, err = op.decodeResult(wm)
		if ep, ok := srvr.(ErrorProcessor); ok {
			ep.ProcessError(err)
		}

		// Pull out $clusterTime and operationTime and update session and clock. We handle this before
		// handling the error to ensure we are properly gossiping the cluster time.
		op.updateClusterTimes(res)
		op.updateOperationTime(res)

		var perr error
		if op.ProcessResponseFn != nil {
			perr = op.ProcessResponseFn(res, srvr)
		}
		switch tt := err.(type) {
		case WriteCommandError:
			if retryable == RetryWrite && tt.Retryable() && retries != 0 {
				retries--
				original, err = err, nil
				conn.Close() // Avoid leaking the connection.
				srvr, err = op.selectServer(ctx)
				if err != nil {
					return original
				}
				conn, err := srvr.Connection(ctx)
				if err != nil || conn == nil || op.retryable(conn.Description()) != RetryWrite {
					if conn != nil {
						conn.Close()
					}
					return original
				}
				defer conn.Close() // Avoid leaking the new connection.
				continue
			}
			// If batching is enabled and either ordered is the default (which is true) or
			// explicitly set to true and we have write errors, return the errors.
			if batching && (op.Batches.Ordered == nil || *op.Batches.Ordered == true) && len(tt.WriteErrors) > 0 {
				return tt
			}
			operationErr.WriteConcernError = tt.WriteConcernError
			operationErr.WriteErrors = append(operationErr.WriteErrors, tt.WriteErrors...)
		case Error:
			if retryable == RetryWrite && tt.Retryable() && retries != 0 {
				retries--
				original, err = err, nil
				conn.Close() // Avoid leaking the connection.
				srvr, err = op.selectServer(ctx)
				if err != nil {
					return original
				}
				conn, err := srvr.Connection(ctx)
				if err != nil || conn == nil || op.retryable(conn.Description()) != RetryWrite {
					if conn != nil {
						conn.Close()
					}
					return original
				}
				defer conn.Close() // Avoid leaking the new connection.
				continue
			}
			return err
		case nil:
			if perr != nil {
				return perr
			}
		default:
			return err
		}

		if batching && len(op.Batches.Documents) > 0 {
			if retryable == RetryWrite && op.Client != nil && op.RetryMode != nil {
				if *op.RetryMode > RetryNone {
					op.Client.IncrementTxnNumber()
				}
				if *op.RetryMode == RetryOncePerCommand {
					retries = 1
				}
			}
			op.Batches.ClearBatch()
			continue
		}
		break
	}
	return nil
}

// Retryable writes are supported if the server supports sessions, the operation is not
// within a transaction, and the write is acknowledged
func (op Operation) retryable(desc description.Server) RetryType {
	switch op.RetryType {
	case RetryWrite:
		if op.Deployment.SupportsRetry() &&
			description.SessionsSupported(desc.WireVersion) &&
			op.Client != nil && !(op.Client.TransactionInProgress() || op.Client.TransactionStarting()) &&
			writeconcern.AckWrite(op.WriteConcern) {
			return RetryWrite
		}
	}
	return RetryType(0)
}

// roundTrip writes a wiremessage to the connection and then reads a wiremessage. The wm parameter
// is reused when reading the wiremessage.
func (op Operation) roundTrip(ctx context.Context, conn Connection, wm []byte) ([]byte, error) {
	err := conn.WriteWireMessage(ctx, wm)
	if err != nil {
		return nil, Error{Message: err.Error(), Labels: []string{TransientTransactionError, NetworkError}}
	}

	res, err := conn.ReadWireMessage(ctx, wm[:0])
	if err != nil {
		err = Error{Message: err.Error(), Labels: []string{TransientTransactionError, NetworkError}}
	}
	return res, err
}

// decompressWireMessage handles decompressing a wiremessage. If the wiremessage
// is not compressed, this method will return the wiremessage.
func (Operation) decompressWireMessage(wm []byte) ([]byte, error) {
	// read the header and ensure this is a compressed wire message
	length, reqid, respto, opcode, rem, ok := wiremessagex.ReadHeader(wm)
	if !ok || len(wm) < int(length) {
		return nil, errors.New("malformed wire message: insufficient bytes")
	}
	if opcode != wiremessage.OpCompressed {
		return wm, nil
	}
	// get the original opcode and uncompressed size
	opcode, rem, ok = wiremessagex.ReadCompressedOriginalOpCode(rem)
	if !ok {
		return nil, errors.New("malformed OP_COMPRESSED: missing original opcode")
	}
	uncompressedSize, rem, ok := wiremessagex.ReadCompressedUncompressedSize(rem)
	if !ok {
		return nil, errors.New("malformed OP_COMPRESSED: missing uncompressed size")
	}
	// get the compressor ID and decompress the message
	compressorID, rem, ok := wiremessagex.ReadCompressedCompressorID(rem)
	if !ok {
		return nil, errors.New("malformed OP_COMPRESSED: missing compressor ID")
	}
	compressedSize := length - 25 // header (16) + original opcode (4) + uncompressed size (4) + compressor ID (1)
	// return the original wiremessage
	msg, rem, ok := wiremessagex.ReadCompressedCompressedMessage(rem, compressedSize)
	if !ok {
		return nil, errors.New("malformed OP_COMPRESSED: insufficient bytes for compressed wiremessage")
	}

	header := make([]byte, 0, uncompressedSize+16)
	header = wiremessagex.AppendHeader(header, uncompressedSize, reqid, respto, opcode)
	uncompressed := make([]byte, uncompressedSize)
	switch compressorID {
	case wiremessage.CompressorSnappy:
		var err error
		uncompressed, err = snappy.Decode(uncompressed, msg)
		if err != nil {
			return nil, err
		}
	case wiremessage.CompressorZLib:
		decompressor, err := zlib.NewReader(bytes.NewReader(msg))
		if err != nil {
			return nil, err
		}
		_, err = io.ReadFull(decompressor, uncompressed)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown compressorID %d", compressorID)
	}
	return append(header, uncompressed...), nil
}

func (op Operation) createWireMessage(dst []byte, desc description.SelectedServer) ([]byte, error) {
	if desc.WireVersion == nil || desc.WireVersion.Max < wiremessage.OpmsgWireVersion {
		return op.createQueryWireMessage(dst, desc)
	}
	return op.createMsgWireMessage(dst, desc)
}

func (op Operation) createQueryWireMessage(dst []byte, desc description.SelectedServer) ([]byte, error) {
	flags := op.slaveOK(desc)
	var wmindex int32
	wmindex, dst = wiremessagex.AppendHeaderStart(dst, wiremessage.NextRequestID(), 0, wiremessage.OpQuery)
	dst = wiremessagex.AppendQueryFlags(dst, flags)
	// FullCollectionName
	dst = append(dst, op.Database...)
	dst = append(dst, dollarCmd[:]...)
	dst = append(dst, 0x00)
	dst = wiremessagex.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessagex.AppendQueryNumberToReturn(dst, -1)

	wrapper := int32(-1)
	rp := op.createReadPref(desc.Server.Kind, desc.Kind, true)
	if len(rp) > 0 {
		wrapper, dst = bsoncore.AppendDocumentStart(dst)
		dst = bsoncore.AppendHeader(dst, bsontype.EmbeddedDocument, "$query")
	}
	idx, dst := bsoncore.AppendDocumentStart(dst)
	dst, err := op.CommandFn(dst, desc)
	if err != nil {
		return dst, err
	}

	if op.Batches != nil && len(op.Batches.Current) > 0 {
		aidx, dst := bsoncore.AppendArrayElementStart(dst, op.Batches.Identifier)
		for i, doc := range op.Batches.Current {
			dst = bsoncore.AppendDocumentElement(dst, strconv.Itoa(i), doc)
		}
		dst, _ = bsoncore.AppendArrayEnd(dst, aidx)
	}

	dst, err = op.addReadConcern(dst, desc)
	if err != nil {
		return dst, err
	}

	dst, err = op.addWriteConcern(dst)
	if err != nil {
		return dst, err
	}

	dst, err = op.addSession(dst, desc)
	if err != nil {
		return dst, err
	}

	// TODO(GODRIVER-617): This should likely be part of addSession, but we need to ensure that we
	// either turn off RetryWrite when we are doing a retryable read or that we pass in RetryType to
	// addSession. We should also only be adding this if the connection supports sessions, but I
	// think that's a given if we've set RetryWrite to true.
	if op.RetryType == RetryWrite && op.Client != nil && op.Client.RetryWrite {
		dst = bsoncore.AppendInt64Element(dst, "txnNumber", op.Client.TxnNumber)
	}

	dst = op.addClusterTime(dst, desc)

	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)

	if len(rp) > 0 {
		var err error
		dst = bsoncore.AppendDocumentElement(dst, "$readPreference", rp)
		dst, err = bsoncore.AppendDocumentEnd(dst, wrapper)
		if err != nil {
			return dst, err
		}
	}

	return bsoncore.UpdateLength(dst, wmindex, int32(len(dst[wmindex:]))), nil
}

func (op Operation) createMsgWireMessage(dst []byte, desc description.SelectedServer) ([]byte, error) {
	// TODO(GODRIVER-617): We need to figure out how to include the writeconcern here so that we can
	// set the moreToCome bit.
	var flags wiremessage.MsgFlag
	var wmindex int32
	wmindex, dst = wiremessagex.AppendHeaderStart(dst, wiremessage.NextRequestID(), 0, wiremessage.OpMsg)
	dst = wiremessagex.AppendMsgFlags(dst, flags)
	// Body
	dst = wiremessagex.AppendMsgSectionType(dst, wiremessage.SingleDocument)

	idx, dst := bsoncore.AppendDocumentStart(dst)

	dst, err := op.CommandFn(dst, desc)
	if err != nil {
		return dst, err
	}
	dst, err = op.addReadConcern(dst, desc)
	if err != nil {
		return dst, err
	}
	dst, err = op.addWriteConcern(dst)
	if err != nil {
		return dst, err
	}

	dst, err = op.addSession(dst, desc)
	if err != nil {
		return dst, err
	}

	// TODO(GODRIVER-617): This should likely be part of addSession, but we need to ensure that we
	// either turn off RetryWrite when we are doing a retryable read or that we pass in RetryType to
	// addSession. We should also only be adding this if the connection supports sessions, but I
	// think that's a given if we've set RetryWrite to true.
	if op.RetryType == RetryWrite && op.Client != nil && op.Client.RetryWrite {
		dst = bsoncore.AppendInt64Element(dst, "txnNumber", op.Client.TxnNumber)
	}

	dst = op.addClusterTime(dst, desc)

	dst = bsoncore.AppendStringElement(dst, "$db", op.Database)
	rp := op.createReadPref(desc.Server.Kind, desc.Kind, false)
	if len(rp) > 0 {
		dst = bsoncore.AppendDocumentElement(dst, "$readPreference", rp)
	}

	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)

	if op.Batches != nil && len(op.Batches.Current) > 0 {
		dst = wiremessagex.AppendMsgSectionType(dst, wiremessage.DocumentSequence)
		idx, dst = bsoncore.ReserveLength(dst)

		dst = append(dst, op.Batches.Identifier...)
		dst = append(dst, 0x00)

		for _, doc := range op.Batches.Current {
			dst = append(dst, doc...)
		}

		dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	}

	return bsoncore.UpdateLength(dst, wmindex, int32(len(dst[wmindex:]))), nil
}

func (op Operation) addReadConcern(dst []byte, desc description.SelectedServer) ([]byte, error) {
	rc := op.ReadConcern
	client := op.Client
	// Starting transaction's read concern overrides all others
	if client != nil && client.TransactionStarting() && client.CurrentRc != nil {
		rc = client.CurrentRc
	}

	// start transaction must append afterclustertime IF causally consistent and operation time exists
	if rc == nil && client != nil && client.TransactionStarting() && client.Consistent && client.OperationTime != nil {
		rc = readconcern.New()
	}

	if rc == nil {
		return dst, nil
	}

	_, data, err := rc.MarshalBSONValue() // always returns a document
	if err != nil {
		return dst, err
	}

	if description.SessionsSupported(desc.WireVersion) && client != nil && client.Consistent && client.OperationTime != nil {
		data = data[:len(data)-1] // remove the null byte
		data = bsoncore.AppendTimestampElement(data, "afterClusterTime", client.OperationTime.T, client.OperationTime.I)
		data, _ = bsoncore.AppendDocumentEnd(data, 0)
	}

	return bsoncore.AppendDocumentElement(dst, "readConcern", data), nil
}

func (op Operation) addWriteConcern(dst []byte) ([]byte, error) {
	wc := op.WriteConcern
	if wc == nil {
		return dst, nil
	}

	t, data, err := wc.MarshalBSONValue()
	if err == writeconcern.ErrEmptyWriteConcern {
		return dst, nil
	}
	if err != nil {
		return dst, err
	}

	return append(bsoncore.AppendHeader(dst, t, "writeConcern"), data...), nil
}

func (op Operation) addSession(dst []byte, desc description.SelectedServer) ([]byte, error) {
	client := op.Client
	if client == nil || !description.SessionsSupported(desc.WireVersion) || desc.SessionTimeoutMinutes == 0 {
		return dst, nil
	}
	if client.Terminated {
		return dst, session.ErrSessionEnded
	}
	lsid, _ := client.SessionID.MarshalBSON()
	dst = bsoncore.AppendDocumentElement(dst, "lsid", lsid)

	if client.TransactionRunning() || client.RetryingCommit {
		dst = bsoncore.AppendInt64Element(dst, "txnNumber", client.TxnNumber)
		if client.TransactionStarting() {
			dst = bsoncore.AppendBooleanElement(dst, "startTransaction", true)
		}
		dst = bsoncore.AppendBooleanElement(dst, "autocommit", false)
	}

	client.ApplyCommand(desc.Server)

	return dst, nil
}

func (op Operation) addClusterTime(dst []byte, desc description.SelectedServer) []byte {
	client, clock := op.Client, op.Clock
	if (clock == nil && client == nil) || !description.SessionsSupported(desc.WireVersion) {
		return dst
	}
	clusterTime := clock.GetClusterTime()
	if client != nil {
		clusterTime = session.MaxClusterTime(clusterTime, client.ClusterTime)
	}
	if clusterTime == nil {
		return dst
	}
	val, err := clusterTime.LookupErr("$clusterTime")
	if err != nil {
		return dst
	}
	return append(bsoncore.AppendHeader(dst, val.Type, "$clusterTime"), val.Value...)
	// return bsoncore.AppendDocumentElement(dst, "$clusterTime", clusterTime)
}

// updateClusterTimes updates the cluster times for the session and cluster clock attached to this
// operation. While the session's AdvanceClusterTime may return an error, this method does not
// because an error being returned from this method will not be returned further up.
func (op Operation) updateClusterTimes(response bsoncore.Document) {
	// Extract cluster time.
	value, err := response.LookupErr("$clusterTime")
	if err != nil {
		// $clusterTime not included by the server
		return
	}
	clusterTime := bsoncore.BuildDocumentFromElements(nil, bsoncore.AppendValueElement(nil, "$clusterTime", value))

	sess, clock := op.Client, op.Clock

	if sess != nil {
		_ = sess.AdvanceClusterTime(bson.Raw(clusterTime))
	}

	if clock != nil {
		clock.AdvanceClusterTime(bson.Raw(clusterTime))
	}
}

// updateOperationTime updates the operation time on the session attached to this operation. While
// the session's AdvanceOperationTime method may return an error, this method does not because an
// error being returned from this method will not be returned further up.
func (op Operation) updateOperationTime(response bsoncore.Document) {
	sess := op.Client
	if sess == nil {
		return
	}

	opTimeElem, err := response.LookupErr("operationTime")
	if err != nil {
		// operationTime not included by the server
		return
	}

	t, i := opTimeElem.Timestamp()
	_ = sess.AdvanceOperationTime(&primitive.Timestamp{
		T: t,
		I: i,
	})
}

func (op Operation) createReadPref(serverKind description.ServerKind, topologyKind description.TopologyKind, isOpQuery bool) bsoncore.Document {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	rp := op.ReadPreference

	if rp == nil {
		if topologyKind == description.Single && serverKind != description.Mongos {
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc
		}
		return nil
	}

	switch rp.Mode() {
	case readpref.PrimaryMode:
		if serverKind == description.Mongos {
			return nil
		}
		if topologyKind == description.Single {
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc
		}
		doc = bsoncore.AppendStringElement(doc, "mode", "primary")
	case readpref.PrimaryPreferredMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
	case readpref.SecondaryPreferredMode:
		_, ok := rp.MaxStaleness()
		if serverKind == description.Mongos && isOpQuery && !ok && len(rp.TagSets()) == 0 {
			return nil
		}
		doc = bsoncore.AppendStringElement(doc, "mode", "secondaryPreferred")
	case readpref.SecondaryMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "secondary")
	case readpref.NearestMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "nearest")
	}

	sets := make([]bsoncore.Document, 0, len(rp.TagSets()))
	for _, ts := range rp.TagSets() {
		if len(ts) == 0 {
			continue
		}
		i, set := bsoncore.AppendDocumentStart(nil)
		for _, t := range ts {
			set = bsoncore.AppendStringElement(set, t.Name, t.Value)
		}
		set, _ = bsoncore.AppendDocumentEnd(set, i)
		sets = append(sets, set)
	}
	if len(sets) > 0 {
		var aidx int32
		aidx, doc = bsoncore.AppendArrayElementStart(doc, "tags")
		for i, set := range sets {
			doc = bsoncore.AppendDocumentElement(doc, strconv.Itoa(i), set)
		}
		doc, _ = bsoncore.AppendArrayEnd(doc, aidx)
	}

	if d, ok := rp.MaxStaleness(); ok {
		doc = bsoncore.AppendInt32Element(doc, "maxStalenessSeconds", int32(d.Seconds()))
	}

	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc
}

func (op Operation) slaveOK(desc description.SelectedServer) wiremessage.QueryFlag {
	if desc.Kind == description.Single && desc.Server.Kind != description.Mongos {
		return wiremessage.SlaveOK
	}

	if rp := op.ReadPreference; rp != nil && rp.Mode() != readpref.PrimaryMode {
		return wiremessage.SlaveOK
	}

	return 0
}

func (Operation) canCompress(cmd string) bool {
	if cmd == "isMaster" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "authenticate" ||
		cmd == "createUser" || cmd == "updateUser" || cmd == "copydbSaslStart" || cmd == "copydbgetnonce" || cmd == "copydb" {
		return false
	}
	return true
}

func (Operation) decodeResult(wm []byte) (bsoncore.Document, error) {
	wmLength := len(wm)
	length, _, _, opcode, wm, ok := wiremessagex.ReadHeader(wm)
	if !ok || int(length) > wmLength {
		return nil, errors.New("malformed wire message: insufficient bytes")
	}

	wm = wm[:wmLength-16] // constrain to just this wiremessage, incase there are multiple in the slice

	switch opcode {
	case wiremessage.OpReply:
		var flags wiremessage.ReplyFlag
		flags, wm, ok = wiremessagex.ReadReplyFlags(wm)
		if !ok {
			return nil, errors.New("malformed OP_REPLY: missing flags")
		}
		_, wm, ok = wiremessagex.ReadReplyCursorID(wm)
		if !ok {
			return nil, errors.New("malformed OP_REPLY: missing cursorID")
		}
		_, wm, ok = wiremessagex.ReadReplyStartingFrom(wm)
		if !ok {
			return nil, errors.New("malformed OP_REPLY: missing startingFrom")
		}
		var numReturned int32
		numReturned, wm, ok = wiremessagex.ReadReplyNumberReturned(wm)
		if !ok {
			return nil, errors.New("malformed OP_REPLY: missing numberReturned")
		}
		if numReturned == 0 {
			return nil, ErrNoDocCommandResponse
		}
		if numReturned > 1 {
			return nil, ErrMultiDocCommandResponse
		}
		var rdr bsoncore.Document
		rdr, rem, ok := wiremessagex.ReadReplyDocument(wm)
		if !ok || len(rem) > 0 {
			return nil, NewCommandResponseError("malformed OP_REPLY: NumberReturned does not match number of documents returned", nil)
		}
		err := rdr.Validate()
		if err != nil {
			return nil, NewCommandResponseError("malformed OP_REPLY: invalid document", err)
		}
		if flags&wiremessage.QueryFailure == wiremessage.QueryFailure {
			return nil, QueryFailureError{
				Message:  "command failure",
				Response: rdr,
			}
		}

		return rdr, extractError(rdr)
	case wiremessage.OpMsg:
		_, wm, ok = wiremessagex.ReadMsgFlags(wm)
		if !ok {
			return nil, errors.New("malformed wire message: missing OP_MSG flags")
		}

		var res bsoncore.Document
		for len(wm) > 0 {
			var stype wiremessage.SectionType
			stype, wm, ok = wiremessagex.ReadMsgSectionType(wm)
			if !ok {
				return nil, errors.New("malformed wire message: insuffienct bytes to read section type")
			}

			switch stype {
			case wiremessage.SingleDocument:
				res, wm, ok = wiremessagex.ReadMsgSectionSingleDocument(wm)
				if !ok {
					return nil, errors.New("malformed wire message: insufficient bytes to read single document")
				}
			case wiremessage.DocumentSequence:
				// TODO(GODRIVER-617): Implement document sequence returns.
				_, _, wm, ok = wiremessagex.ReadMsgSectionDocumentSequence(wm)
				if !ok {
					return nil, errors.New("malformed wire message: insufficient bytes to read document sequence")
				}
			default:
				return nil, fmt.Errorf("malformed wire message: uknown section type %v", stype)
			}
		}

		err := res.Validate()
		if err != nil {
			return nil, NewCommandResponseError("malformed OP_MSG: invalid document", err)
		}

		return res, extractError(res)
	default:
		return nil, fmt.Errorf("cannot decode result from %s", opcode)
	}
}
