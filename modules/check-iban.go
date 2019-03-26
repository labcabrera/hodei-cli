package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const CheckIbanCmd = "check-iban"

type CheckIbanOptions struct {
	CountryCode string
	Iban        string
	Help        bool
	Verbose     bool
}

func CheckIban(options *CheckIbanOptions) (res string, err error) {
	if options.Verbose {
		log.Printf("Validating IBAN %s", options.Iban)
	}
	headers := amqp.Table{
		"App-Source": "hodei-cli",
	}
	body := `{"countryCode": "` + options.CountryCode + `","iban": "` + options.Iban + `"}`
	res, err = client.SendAndReceive("cnp.sepa", "iban.validation", body, headers, options.Verbose)
	return
}

func CheckIbanFlagSet(options *CheckIbanOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(CheckIbanCmd, flag.ExitOnError)
	fs.StringVar(&options.Iban, "iban", "", "IBAN")
	fs.StringVar(&options.CountryCode, "country", "", "Country ISO3 code")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
