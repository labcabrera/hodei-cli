package modules

import (
	"flag"
	"fmt"
	"os"

	"github.com/labcabrera/hodei-cli/client"
	"github.com/streadway/amqp"
)

const PullPoliciesCmd = "pull-policies"

type PullPoliciesOptions struct {
	Product      string
	Id           string
	ExternalCode string
	AgreementId  string
	Username     string
	Authorities  string
	Verbose      bool
	Help         bool
}

func PullPolicies(options *PullPoliciesOptions) {
	if options.Product == "" {
		fmt.Println("Missing product parameter")
		os.Exit(1)
	} else if options.Username == "" || options.Authorities == "" {
		fmt.Println("Missing security parameters")
		os.Exit(1)
	}
	productMapping := map[string]string{
		"ppi": "ppi.referential",
	}
	exchange := productMapping[options.Product]

	body := `{"id": "` + options.Id + `", "externalCode": "` + options.ExternalCode + `", "agreementId":"` + options.AgreementId + `"}`
	headers := amqp.Table{
		"App-Username":    options.Username,
		"App-Authorities": options.Authorities,
	}
	client.SendMessageWithHeaders(exchange, "policy.pull", body, headers, options.Verbose)
	return
}

func PullPoliciesFlagSet(options *PullPoliciesOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(PullPoliciesCmd, flag.ExitOnError)
	fs.StringVar(&options.Product, "product", "", "Product external code")
	fs.StringVar(&options.Id, "id", "", "Policy identifier")
	fs.StringVar(&options.ExternalCode, "externalcode", "", "Policy external code")
	fs.StringVar(&options.AgreementId, "agreement", "", "Agreement identifier")
	fs.StringVar(&options.Username, "u", "", "Username")
	fs.StringVar(&options.Authorities, "a", "", "Authorities")
	fs.BoolVar(&options.Verbose, "v", false, "Verbose")
	fs.BoolVar(&options.Help, "help", false, "Help")
	return fs
}
