package main

import (
	"fmt"
	"flag"
	"log"
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/labcabrera/hodei-cli/modules"
	"github.com/labcabrera/hodei-cli/client"
)

const version = "0.2.0-SNAPSHOT"
const exchangeReferential = "cnp.referential"

type PolicyPullRequest struct {
	entityId			string
	externalCode		string
	agreementId			string
}

type Authorization struct {
	username			string
	authorities			string
}

func main() {

	// Arguments
	pullCountries  := flag.Bool("pull-countries",  false, "Pull countries over referential API")
	pullProducts   := flag.Bool("pull-products",   false, "Pull products and agreements over referential API")
	pullPolicies   := flag.Bool("pull-policies",   false, "Pull policies")
	pullPerson     := flag.Bool("pull-person",     false, "Pull person from referential API")
	pullLegal      := flag.Bool("pull-legal",      false, "Pull legal entity from referential API")

	id             := flag.String("id",            "",    "Entity MongoDB identifier")
	externalCode   := flag.String("external-code", "",    "Entity external code")
	product        := flag.String("product",       "",    "Product name")
	agreement      := flag.String("agreement",     "",    "Agreement external code")
	argIban           := flag.String("iban",          "",    "IBAN validation")

	username       := flag.String("u",             "",    "Username")
	authorities    := flag.String("a",             "",    "Authorities (coma separated list)")
	verbose        := flag.Bool("v",               false, "Verbose")
	printVersion   := flag.Bool("version",         false, "Print version")
	
	flag.Parse()
	
	if(*printVersion) {
		fmt.Println("Hodei cli ", version)
		return
	}
	if(*argIban != "") {
		iban.ProcessIban(*argIban, *verbose)
		return
	}
	if(*pullCountries) {
		sendPullCountryMessage(*verbose)
	}
	if(*pullProducts) {
		sendPullProductsMessage(*verbose)
	}
	if(*pullPerson) {
		//TODO
		fmt.Println("Not implemented")
	}
	if(*pullLegal) {
		//TODO
		fmt.Println("Not implemented")
	}
	if(*pullPolicies) {
		request := PolicyPullRequest{*id, *externalCode, *agreement}
		auth := Authorization{*username, *authorities}
		sendPullPoliciesMessage(*product, request, auth, *verbose)
	}
}

func sendPullCountryMessage(verbose bool) {
	if(verbose) {
		log.Printf("Pulling countries from referential API")
	}
	client.SendMessage(exchangeReferential, "country.pull", "", verbose)
}

func sendPullProductsMessage(verbose bool) {
	if(verbose) {
		log.Printf("Pulling products and agreements from referential API")
	}
	client.SendMessage(exchangeReferential, "product.pull", "", verbose)
}

func sendPullPoliciesMessage(product string, request PolicyPullRequest, auth Authorization, verbose bool) {
	switch(product) {
	case "":
		fmt.Println("Required argument product")
		return
	case "ppi":
		log.Printf("Agreement: %s", request.agreementId)

		bodyBinary, err := json.Marshal(request)
		if(err != nil) {
			log.Fatalf("%s: %s", "Error marshalling request", err)
			return
		}
		body := string(bodyBinary)

		//TODO
		body = `{"agreementId":"` + request.agreementId + `"}`		
		headers := amqp.Table{
			"App-Username"   : auth.username,
			"App-Authorities": auth.authorities,
		}
		client.SendMessageWithHeaders("ppi.referential", "policy.pull", body, headers, verbose)
	default:
		log.Fatalf("Unknown product %s", product)
	}
}
