package modules

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

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
	switch options.Product {
	case "":
		fmt.Println("Required argument product")
		return
	case "ppi":
		log.Printf("Agreement: %s", options.AgreementId)

		bodyBinary, err := json.Marshal(options)
		if err != nil {
			log.Fatalf("%s: %s", "Error marshalling request", err)
			return
		}
		body := string(bodyBinary)

		//TODO
		body = `{"agreementId":"` + options.AgreementId + `"}`
		headers := amqp.Table{
			"App-Username":    options.Username,
			"App-Authorities": options.Authorities,
		}
		client.SendMessageWithHeaders("ppi.referential", "policy.pull", body, headers, options.Verbose)
	default:
		log.Fatalf("Unknown product %s", options.Product)
	}
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
