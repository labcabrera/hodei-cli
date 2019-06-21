package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const PullCoveragesCmd = "pull-coverages"

type PullCoveragesOptions struct {
	Id                 string
	ExternalCode       string
	PolicyId           string
	PolicyExternalCode string
	Username           string
	Authorities        string
	Verbose            bool
	Help               bool
}

func PullCoverages(options *PullCoveragesOptions) (res string, err error) {
	if options.Verbose {
		log.Printf("Pulling coverages from referential API")
	}
	headers := amqp.Table{
		"App-Username":    options.Username,
		"App-Authorities": options.Authorities,
	}
	body := `{"id":"` + options.Id +
		`","externalCode":"` + options.ExternalCode +
		`","policyId":"` + options.PolicyId +
		`","policyExternalCode":"` + options.PolicyExternalCode +
		`"}`
	client.SendMessageWithHeaders("cnp.referential", "coverage.pull", body, headers, options.Verbose)
	return
}

func PullCoveragesFlagSet(options *PullCoveragesOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullCoveragesCmd, flag.ExitOnError)
	fs.StringVar(&options.Id, "id", "", "Coverage identifier")
	fs.StringVar(&options.ExternalCode, "externalcode", "", "Coverage external code")
	fs.StringVar(&options.PolicyId, "policyid", "", "Policy identifier")
	fs.StringVar(&options.PolicyExternalCode, "policyexternalcode", "", "Policy external code")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	return fs
}
