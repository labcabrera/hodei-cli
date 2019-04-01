package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const PullProductsCmd = "pull-products"

type PullProductsOptions struct {
	Username    string
	Authorities string
	Verbose     bool
	Help        bool
}

func PullProducts(options *PullProductsOptions) {
	if options.Verbose {
		log.Printf("Pulling products from referential API")
	}
	headers := amqp.Table{
		"App-Username":    options.Username,
		"App-Authorities": options.Authorities,
	}
	body := "{}"
	client.SendMessageWithHeaders("cnp.referential", "product.pull", body, headers, options.Verbose)
}

func PullProductsFlagSet(options *PullProductsOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullProductsCmd, flag.ExitOnError)
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	return fs
}
