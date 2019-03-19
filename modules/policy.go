package modules

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/labcabrera/hodei-cli/client"
)

type PolicyPullRequest struct {
	EntityId			string
	ExternalCode		string
	AgreementId			string
}

type Authorization struct {
	Username			string
	Authorities			string
}

func PullPolicies(product string, request PolicyPullRequest, auth Authorization, verbose bool) {
	switch(product) {
	case "":
		fmt.Println("Required argument product")
		return
	case "ppi":
		log.Printf("Agreement: %s", request.AgreementId)

		bodyBinary, err := json.Marshal(request)
		if(err != nil) {
			log.Fatalf("%s: %s", "Error marshalling request", err)
			return
		}
		body := string(bodyBinary)

		//TODO
		body = `{"agreementId":"` + request.AgreementId + `"}`		
		headers := amqp.Table{
			"App-Username"   : auth.Username,
			"App-Authorities": auth.Authorities,
		}
		client.SendMessageWithHeaders("ppi.referential", "policy.pull", body, headers, verbose)
	default:
		log.Fatalf("Unknown product %s", product)
	}
}

