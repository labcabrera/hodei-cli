package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

type PullAgreementsOptions struct {
	Id           string
	ExternalCode string
	Product      string
	Username     string
	Authorities  string
	Verbose      bool
	Help         bool
}

const PullAgreementsCmd = "pull-agreements"

func PullAgreements(options *PullAgreementsOptions) {
	if options.Verbose {
		log.Printf("Pulling agreements from referential API")
	}
	headers := amqp.Table{
		"App-Username":    options.Username,
		"App-Authorities": options.Authorities,
	}
	body := "{}"
	client.SendMessageWithHeaders("cnp.referential", "agreement.pull", body, headers, options.Verbose)
}

func PullAgreementsFlagSet(options *PullAgreementsOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullAgreementsCmd, flag.ExitOnError)
	fs.StringVar(&options.Id, "id", "", "Agreement identifier")
	fs.StringVar(&options.ExternalCode, "externalcode", "", "Agreement external code")
	fs.StringVar(&options.ExternalCode, "product", "", "Product identifier")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
