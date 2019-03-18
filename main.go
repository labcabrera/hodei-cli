package main

import (
	"fmt"
	"flag"
	"log"
	"github.com/labcabrera/hodei-cli/modules"
	"github.com/labcabrera/hodei-cli/client"
)

const version = "0.2.0-SNAPSHOT"
const exchangeReferential = "cnp.referential"

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
	argIban           := flag.String("iban",       "",    "IBAN validation")

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
		modules.CheckIban(*argIban, *verbose)
		return
	}
	if(*pullCountries) {
		modules.PullCountries(*verbose)
	}
	if(*pullProducts) {
		modules.PullProducts(*verbose)
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
		request := modules.PolicyPullRequest{*id, *externalCode, *agreement}
		auth := modules.Authorization{*username, *authorities}
		modules.PullPolicies(*product, request, auth, *verbose)
	}
}

func sendPullCountryMessage(verbose bool) {
	if(verbose) {
		log.Printf("Pulling countries from referential API")
	}
	client.SendMessage(exchangeReferential, "country.pull", "", verbose)
}
