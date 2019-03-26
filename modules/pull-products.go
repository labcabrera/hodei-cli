package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
)

const PullProductsCmd = "pull-products"

type PullProductsOptions struct {
	Verbose bool
	Help    bool
}

func PullProducts(options *PullProductsOptions) {
	if options.Verbose {
		log.Printf("Pulling products from referential API")
	}
	client.SendMessage("cnp.referential", "product.pull", "", options.Verbose)
}

func PullProductsFlagSet(options *PullProductsOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullProductsCmd, flag.ExitOnError)
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
