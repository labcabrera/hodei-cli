package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
)

const PullProfessionsCmd = "pull-professions"

type PullProfessionsOptions struct {
	Help    bool
	Verbose bool
}

func PullProfessions(options *PullProfessionsOptions) {
	if options.Verbose {
		log.Printf("Pulling professions from referential API")
	}
	client.SendMessage("cnp.referential", "profession.pull", "{}", options.Verbose)
}

func PullProfessionsFlagSet(options *PullProfessionsOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullProfessionsCmd, flag.ExitOnError)
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
