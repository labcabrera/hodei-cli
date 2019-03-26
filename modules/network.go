package modules

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

type PullNetworkOptions struct {
	Id           string
	ExternalCode string
	IdCard       string
	Username     string
	Authorities  string
	Verbose      bool
	Help         bool
}

func PullNetworks(options *PullNetworkOptions) {
	if options.Verbose {
		log.Printf("Pulling networks")
		log.Printf("Username: %s", options.Username)
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
	client.SendMessageWithHeaders("cnp.referential", "network.pull", body, headers, options.Verbose)
	return
}

func NetworkFlagset(options *PullNetworkOptions) *flag.FlagSet {
	fs := flag.NewFlagSet("pull-networks", flag.ExitOnError)
	fs.StringVar(&options.Id, "id", "", "Entity identifier")
	fs.StringVar(&options.ExternalCode, "externalcode", "", "Entity external code")
	fs.StringVar(&options.IdCard, "idcard", "", "Entity IdCard")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
