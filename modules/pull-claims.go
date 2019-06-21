package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const PullClaimsCmd = "pull-claims"

type PullClaimsOptions struct {
	Id                 string
	ExternalCode       string
	PolicyId           string
	PolicyExternalCode string
	Username           string
	Authorities        string
	Verbose            bool
	Help               bool
}

func PullClaims(options *PullClaimsOptions) (res string, err error) {
	if options.Verbose {
		log.Printf("Pulling claims from referential API")
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
	client.SendMessageWithHeaders("cnp.referential", "claim.pull", body, headers, options.Verbose)
	return
}

func PullClaimsFlagSet(options *PullClaimsOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullClaimsCmd, flag.ExitOnError)
	fs.StringVar(&options.Id, "id", "", "Claim identifier")
	fs.StringVar(&options.ExternalCode, "externalcode", "", "Claim external code")
	fs.StringVar(&options.PolicyId, "policyid", "", "Policy identifier")
	fs.StringVar(&options.PolicyExternalCode, "policyexternalcode", "", "Policy external code")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	return fs
}
