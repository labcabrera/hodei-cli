package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
)

const PullCountriesCmd = "pull-countries"

type PullCountriesOptions struct {
	Verbose bool
	Help    bool
}

func PullCountries(options *PullCountriesOptions) {
	if options.Verbose {
		log.Printf("Pulling countries from referential API")
	}
	client.SendMessage("cnp.referential", "country.pull", "", options.Verbose)
}

func PullCountriesFlagSet(options *PullCountriesOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullCountriesCmd, flag.ExitOnError)
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
