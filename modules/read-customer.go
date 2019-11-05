package modules

import (
	"flag"
	"log"
	"os"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const CustomerSearchCmd = "read-customer"

type CustomerSearchModule struct {
}

type customerSearchOptions struct {
	id          string
	legal       bool
	username    string
	authorities string
	verbose     bool
	help        bool
}

func (m CustomerSearchModule) Execute(args []string) {
	options := customerSearchOptions{}
	flagset := customerSearchCreateFlagSet(&options)
	flagset.Parse(args)
	if flagset.Parsed() {
		if options.help {
			flagset.PrintDefaults()
		} else {
			customerSearch(&options)
		}
	}
}

func customerSearch(options *customerSearchOptions) (res string, err error) {
	if options.verbose {
		log.Printf("Searching customer %s (%s:%s)", options.id, options.username, options.authorities)
	}
	if options.id == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	personType := "person"
	if options.legal {
		personType = "legal"
	}
	headers := amqp.Table{
		"App-Username":    options.username,
		"App-Authorities": options.authorities,
	}
	body := `{"1":{"type":"` + personType + `","reference":"` + options.id + `"}}`
	res, err = client.SendAndReceive("cnp.customer", "customer.search", body, headers, options.verbose)
	if err != nil {
		log.Fatalf("%s: %s", "Error reading person", err)
		return
	}
	return
}

func customerSearchCreateFlagSet(options *customerSearchOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(CustomerSearchCmd, flag.ExitOnError)
	fs.StringVar(&options.id, "id", "", "Entity identifier")
	fs.StringVar(&options.username, "u", "", "Username")
	fs.StringVar(&options.authorities, "a", "", "Authorities")
	fs.BoolVar(&options.verbose, "v", false, "Verbose")
	fs.BoolVar(&options.help, "help", false, "Help")
	return fs
}
