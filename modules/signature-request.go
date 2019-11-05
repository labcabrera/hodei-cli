package modules

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const SignatureRequestCmd = "signature-request"

type SignatureRequestModule struct {
}

type signatureRequestOptions struct {
	documentId  string
	username    string
	authorities string
	verbose     bool
	help        bool
}

func (m SignatureRequestModule) Execute(args []string) {
	options := signatureRequestOptions{}
	flagset := signatureRequestModuleCreateFlagSet(&options)
	flagset.Parse(args)
	if flagset.Parsed() {
		if options.help {
			flagset.PrintDefaults()
		} else {
			signatureRequest(&options)
		}
	}
}

func signatureRequest(options *signatureRequestOptions) (res string, err error) {
	if options.verbose {
		log.Printf("Sending signature request")
	}
	if options.documentId == "" {
		fmt.Println("Required document identifier")
		os.Exit(1)
	} else if options.username == "" || options.authorities == "" {
		fmt.Println("Required authorization information")
		os.Exit(1)
	}
	headers := amqp.Table{
		"App-Username":    options.username,
		"App-Authorities": options.authorities,
	}
	body := `{"documentId":"` + options.documentId + `"}`
	res, err = client.SendAndReceive("cnp.esignature", "signature.request", body, headers, options.verbose)
	return
}

func signatureRequestModuleCreateFlagSet(options *signatureRequestOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(SignatureRequestCmd, flag.ExitOnError)
	fs.StringVar(&options.documentId, "id", "", "Document identifier")
	fs.StringVar(&options.username, "u", "", "Username")
	fs.StringVar(&options.authorities, "a", "", "Authorities")
	fs.BoolVar(&options.verbose, "v", false, "Verbose")
	fs.BoolVar(&options.help, "help", false, "Help")
	return fs
}
