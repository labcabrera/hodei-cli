package modules

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

type SignatureRequestOptions struct {
	DocumentId  string
	Username    string
	Authorities string
	Verbose     bool
	Help        bool
}

const SignatureRequestCmd = "signature-request"

func SignatureRequest(options *SignatureRequestOptions) (res string, err error) {
	if options.Verbose {
		log.Printf("Sending signature request")
	}
	if options.DocumentId == "" {
		fmt.Println("Required document identifier")
		os.Exit(1)
	} else if options.Username == "" || options.Authorities == "" {
		fmt.Println("Required authorization information")
		os.Exit(1)
	}
	headers := amqp.Table{
		"App-Username":    options.Username,
		"App-Authorities": options.Authorities,
	}
	body := `{"documentId":"` + options.DocumentId + `"}`
	res, err = client.SendAndReceive("cnp.esignature", "signature.request", body, headers, options.Verbose)
	return
}

func SignatureRequestFlagSet(options *SignatureRequestOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullAgreementsCmd, flag.ExitOnError)
	fs.StringVar(&options.DocumentId, "id", "", "Document identifier")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
