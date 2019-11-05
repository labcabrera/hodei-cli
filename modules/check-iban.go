package modules

import (
	"flag"
	"fmt"
	"log"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const CheckIbanCmd = "check-iban"

type CheckIbanModule struct {
}

type checkIbanOptions struct {
	countryCode string
	iban        string
	help        bool
	verbose     bool
}

func (m CheckIbanModule) Execute(args []string) {
	options := checkIbanOptions{}
	flagset := checkIbanCreateFlagSet(&options)
	flagset.Parse(args)
	if flagset.Parsed() {
		if options.help {
			flagset.PrintDefaults()
		} else {
			checkIban(&options)
		}
	}
}

func checkIban(options *checkIbanOptions) (res string, err error) {
	if options.verbose {
		log.Printf("Validating IBAN %s", options.iban)
	}
	headers := amqp.Table{
		"App-Source": "hodei-cli",
	}
	body := `{"countryCode": "` + options.countryCode + `","iban": "` + options.iban + `"}`
	res, err = client.SendAndReceive("cnp.sepa", "iban.validation", body, headers, options.verbose)
	fmt.Println(res)
	return
}

func checkIbanCreateFlagSet(options *checkIbanOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(CheckIbanCmd, flag.ExitOnError)
	fs.StringVar(&options.iban, "iban", "", "IBAN")
	fs.StringVar(&options.countryCode, "country", "", "Country ISO3 code")
	fs.BoolVar(&options.verbose, "v", false, "Verbose")
	fs.BoolVar(&options.help, "help", false, "Help")
	return fs
}
