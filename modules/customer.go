package modules

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

type PullCustomerOptions struct {
	Id           string
	ExternalCode string
	IdCard       string
	Username     string
	Authorities  string
	Verbose      bool
	Help         bool
}

func PullCustomers(options *PullCustomerOptions) {
	if options.Verbose {
		log.Printf("Pulling customers")
	}
	if options.Id == "" && options.ExternalCode == "" && options.IdCard == "" {
		fmt.Println("Required one pull search method parameter")
		os.Exit(1)
	} else if options.Username == "" || options.Authorities == "" {
		fmt.Println("Required authentication arguments")
		os.Exit(1)
	}
	headers := amqp.Table{
		"App-Username":    options.Username,
		"App-Authorities": options.Authorities,
	}
	body := `{"id": "` + options.Id + `","externalCode": "` + options.ExternalCode + `","idCard": "` + options.IdCard + `"}`
	client.SendMessageWithHeaders("cnp.referential", "customer.pull", body, headers, options.Verbose)
	return
}

func PullCustomersFlagSet(options *PullCustomerOptions) *flag.FlagSet {
	fs := flag.NewFlagSet("pull-customers", flag.ExitOnError)
	fs.StringVar(&options.Id, "id", "", "Entity identifier")
	fs.StringVar(&options.ExternalCode, "externalcode", "", "Entity external code")
	fs.StringVar(&options.IdCard, "idcard", "", "Entity IdCard")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}

func CustomerSearch(id string, legal bool, username string, authorities string, verbose bool) (res string, err error) {
	if verbose {
		log.Printf("Searching customer %s (%s:%s)", id, username, authorities)
	}
	if id == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	personType := "person"
	if legal {
		personType = "legal"
	}

	headers := amqp.Table{
		"App-Username":    username,
		"App-Authorities": authorities,
	}
	body := `{"1":{"type":"` + personType + `","reference":"` + id + `"}}`

	if verbose {
		log.Printf("Body: %s", body)
	}

	res, err = client.SendAndReceive("cnp.customer", "customer.search", body, headers, verbose)
	if err != nil {
		log.Fatalf("%s: %s", "Error reading person", err)
		return
	}

	return
}
