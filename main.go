package main

import (
	"fmt"
	"os"
	"flag"
	//"log"
	"time"
	"math/rand"
	"github.com/labcabrera/hodei-cli/modules"
)

const version = "0.2.0-SNAPSHOT"

func main() {

	readPersonCommand			:= flag.NewFlagSet("read-person", flag.ExitOnError)
	checkIbanCommand            := flag.NewFlagSet("check-iban", flag.ExitOnError)
	pullCountriesCommand		:= flag.NewFlagSet("pull-countries", flag.ExitOnError)

	readPersonIdPtr				:= readPersonCommand.String(	"id",		"", 	"Person MongoDB identifier (required)")
	readPersonUsernamePtr		:= readPersonCommand.String(	"u",		"", 	"Username (required)")
	readPersonAuthoritiesPtr	:= readPersonCommand.String(	"a",		"",		"Authorities (required)")
	readPersonVerbosePtr		:= readPersonCommand.Bool(		"v",		false,	"Verbose")

	checkIbanValuePtr			:= checkIbanCommand.String(		"iban",		"", 	"IBAN (required)")
	checkIbanCountryPtr			:= checkIbanCommand.String(		"c",		"ESP", 	"Country")
	checkIbanVerbosePtr			:= checkIbanCommand.Bool(		"v",		false,	"Verbose")

	pullCountriesVerbosePtr		:= pullCountriesCommand.Bool(	"v",		false,	"Verbose")

	if len(os.Args) < 2 {
		usage()
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())

	switch os.Args[1] {
	case "version":
		fmt.Println("Hodei cli", version)
		os.Exit(0)
	case "read-person":
		readPersonCommand.Parse(os.Args[2:])
	case "check-iban":
		checkIbanCommand.Parse(os.Args[2:])
	case "pull-countries":
		pullCountriesCommand.Parse(os.Args[2:]) 
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	if readPersonCommand.Parsed() {
		if *readPersonIdPtr == "" || *readPersonUsernamePtr == "" || *readPersonAuthoritiesPtr == "" {
			readPersonCommand.PrintDefaults()
			os.Exit(1)
		}
		modules.CustomerSearch(
			*readPersonIdPtr,
			*readPersonUsernamePtr,
			*readPersonAuthoritiesPtr,
			*readPersonVerbosePtr)
	}

	if checkIbanCommand.Parsed() {
		if *checkIbanValuePtr == "" {
			checkIbanCommand.PrintDefaults()
			os.Exit(1)
		}
		modules.CheckIban(*checkIbanCountryPtr, *checkIbanValuePtr, *checkIbanVerbosePtr)
	}

	if pullCountriesCommand.Parsed() {
		modules.PullCountries(*pullCountriesVerbosePtr)
	}


	/*

	// Arguments
	pullCountries	:= flag.Bool("pull-countries",	false,	"Pull countries over referential API")
	pullProducts	:= flag.Bool("pull-products",	false,	"Pull products and agreements over referential API")
	pullPolicies	:= flag.Bool("pull-policies",	false,	"Pull policies")
	pullPerson		:= flag.Bool("pull-person",		false,	"Pull person from referential API")
	pullLegal		:= flag.Bool("pull-legal",		false,	"Pull legal entity from referential API")

	
	id				:= flag.String("id",			"",		"Entity MongoDB identifier")
	externalCode	:= flag.String("external-code",	"",		"Entity external code")
	product			:= flag.String("product",		"",		"Product name")
	agreement		:= flag.String("agreement",		"",		"Agreement external code")
	argIban			:= flag.String("iban",			"",		"IBAN validation")

	username		:= flag.String("u",				"",		"Username")
	authorities		:= flag.String("a",				"",		"Authorities (coma separated list)")
	verbose			:= flag.Bool("v",				false,	"Verbose")
	printVersion	:= flag.Bool("version",			false,	"Print version")
	
	flag.Parse()
	
	if(*printVersion) {
		fmt.Println("Hodei cli ", version)
		return
	}
	if(*argIban != "") {
		msg, err := modules.CheckIban(*argIban, *verbose)
		if(err != nil) {
			log.Fatalf("%s: %s", "IBAN error", err)
		} else {
			fmt.Println(msg)
		}
		return
	}
	if(*pullCountries) {
		modules.PullCountries(*verbose)
	}
	if(*pullProducts) {
		modules.PullProducts(*verbose)
	}
	if(*pullPerson || *pullLegal) {
		fmt.Println("Hodei cli ", version)
	}
	if(*pullPolicies) {
		request := modules.PolicyPullRequest{*id, *externalCode, *agreement}
		auth := modules.Authorization{*username, *authorities}
		modules.PullPolicies(*product, request, auth, *verbose)
	}
	*/
}

/*
func sendPullCountryMessage(verbose bool) {
	if(verbose) {
		log.Printf("Pulling countries from referential API")
	}
	client.SendMessage(exchangeReferential, "country.pull", "", verbose)
}
*/

func usage() {
	fmt.Println(`
Usage: hodei-cli COMMAND [OPTIONS]")

Commands:
  read-person
  read-legal
  pull-countries
  pull-products
  pull-policies
  check-iban
  version
`)
}

