package modules

import (
	"flag"
	"log"
	"os"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const CustomerSearchCmd = "read-customer"

type CustomerSearchOptions struct {
	Id          string
	Legal       bool
	Username    string
	Authorities string
	Verbose     bool
	Help        bool
}

func CustomerSearch(options *CustomerSearchOptions) (res string, err error) {
	if options.Verbose {
		log.Printf("Searching customer %s (%s:%s)", options.Id, options.Username, options.Authorities)
	}
	if options.Id == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	personType := "person"
	if options.Legal {
		personType = "legal"
	}
	headers := amqp.Table{
		"App-Username":    options.Username,
		"App-Authorities": options.Authorities,
	}
	body := `{"1":{"type":"` + personType + `","reference":"` + options.Id + `"}}`
	res, err = client.SendAndReceive("cnp.customer", "customer.search", body, headers, options.Verbose)
	if err != nil {
		log.Fatalf("%s: %s", "Error reading person", err)
		return
	}
	return
}

func CustomerSearchFlagSet(options *CustomerSearchOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(CustomerSearchCmd, flag.ExitOnError)
	fs.StringVar(&options.Id, "id", "", "Entity identifier")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
