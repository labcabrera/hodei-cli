package modules

import (
	"flag"
	"log"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const PullOrdersCmd = "pull-orders"

type PullOrdersOptions struct {
	Id                 string
	ExternalCode       string
	PolicyId           string
	PolicyExternalCode string
	Username           string
	Authorities        string
	Verbose            bool
	Help               bool
}

func PullOrders(options *PullOrdersOptions) {
	if options.Verbose {
		log.Printf("Pulling orders from referential API")
	}
	headers := amqp.Table{
		"App-Username":    options.Username,
		"App-Authorities": options.Authorities,
	}
	body := `{"id": "` + options.Id +
		`","externalCode": "` + options.ExternalCode +
		`","policyId":"` + options.PolicyId +
		`","policyExternalCode":"` + options.PolicyExternalCode +
		`"}`
	client.SendMessageWithHeaders("cnp.referential", "order.pull", body, headers, options.Verbose)
}

func PullOrdersFlagSet(options *PullOrdersOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullOrdersCmd, flag.ExitOnError)
	fs.StringVar(&options.Id, "id", "", "Order identifier")
	fs.StringVar(&options.ExternalCode, "externalcode", "", "Order external code")
	fs.StringVar(&options.PolicyId, "policyid", "", "Policy identifier")
	fs.StringVar(&options.PolicyExternalCode, "policyexternalcode", "", "Policy external code")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	return fs
}
