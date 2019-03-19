package main

import (
	"fmt"
	"os"
	"flag"
	"time"
	"math/rand"
	"github.com/labcabrera/hodei-cli/modules"
)

const version = "0.3.0-SNAPSHOT"

func main() {

	readPersonCommand			:= flag.NewFlagSet("read-person", flag.ExitOnError)
	pullCountriesCommand		:= flag.NewFlagSet("pull-countries", flag.ExitOnError)
	pullProductsCommand			:= flag.NewFlagSet("pull-products", flag.ExitOnError)
	pullPoliciesCommand			:= flag.NewFlagSet("pull-policies", flag.ExitOnError)
	checkIbanCommand			:= flag.NewFlagSet("check-iban", flag.ExitOnError)

	readPersonIdPtr				:= readPersonCommand.String(	"id",		"", 	"MongoDB ID (required)")
	readPersonUsernamePtr		:= readPersonCommand.String(	"u",		"", 	"Username (required)")
	readPersonAuthoritiesPtr	:= readPersonCommand.String(	"a",		"",		"Authorities (required)")
	readPersonLegalPtr			:= readPersonCommand.Bool(		"legal",	false,	"Legal entity")
	readPersonVerbosePtr		:= readPersonCommand.Bool(		"v",		false,	"Verbose")
	readPersonHelpPtr			:= readPersonCommand.Bool(		"help",		false,	"Help")
	
	pullCountriesVerbosePtr		:= pullCountriesCommand.Bool(	"v",		false,	"Verbose")
	pullCountriesHelpPtr		:= pullCountriesCommand.Bool(	"help",		false,	"Help")

	pullProductsVerbosePtr		:= pullPoliciesCommand.Bool(	"v",		false,	"Verbose")
	pullProductsHelpPtr			:= pullPoliciesCommand.Bool(	"help",		false,	"Help")
	
	pullPoliciesProduct			:= pullProductsCommand.String(	"product",	"",		"Product ID (required)")
	pullPoliciesUsernamePtr		:= pullProductsCommand.String(	"u",		"",		"Username (required)")
	pullPoliciesAuthoritiesPtr	:= pullProductsCommand.String(	"a",		"",		"Authorities (required)")
	pullPoliciesIdPtr			:= pullProductsCommand.String(	"id",		"",		"Policy ID")
	pullPoliciesExternalCodePtr	:= pullProductsCommand.String(	"code",		"",		"Policy external code")
	pullPoliciesAgremmentPtr	:= pullProductsCommand.String(	"agreement","",		"Agreement ID")
	pullPoliciesVerbosePtr		:= pullProductsCommand.Bool(	"v",		false,	"Verbose")
	pullPoliciesHelpPtr			:= pullProductsCommand.Bool(	"help",		false,	"Help")
	
	checkIbanValuePtr			:= checkIbanCommand.String(		"iban",		"", 	"IBAN (required)")
	checkIbanCountryPtr			:= checkIbanCommand.String(		"c",		"ESP", 	"Country")
	checkIbanVerbosePtr			:= checkIbanCommand.Bool(		"v",		false,	"Verbose")
	checkIbanHelpPtr			:= checkIbanCommand.Bool(		"help",		false,	"Help")
	
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
	case "pull-products":
		pullProductsCommand.Parse(os.Args[2:])
	case "pull-policies":
		pullPoliciesCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	if readPersonCommand.Parsed() {
		if *readPersonHelpPtr {
			readPersonCommand.PrintDefaults()
			os.Exit(0)
		} else if *readPersonIdPtr == "" || *readPersonUsernamePtr == "" || *readPersonAuthoritiesPtr == "" {
			readPersonCommand.PrintDefaults()
			os.Exit(1)
		}
		modules.CustomerSearch(
			*readPersonIdPtr,
			*readPersonLegalPtr,
			*readPersonUsernamePtr,
			*readPersonAuthoritiesPtr,
			*readPersonVerbosePtr)
	}

	if pullCountriesCommand.Parsed() {
		if *pullCountriesHelpPtr {
			pullCountriesCommand.PrintDefaults()
			os.Exit(0)
		}
		modules.PullCountries(*pullCountriesVerbosePtr)
	}

	if pullProductsCommand.Parsed() {
		if *pullProductsHelpPtr {
			pullProductsCommand.PrintDefaults()
			os.Exit(0)
		}
		modules.PullProducts(*pullProductsVerbosePtr)
	}
	
	if pullPoliciesCommand.Parsed() {
		if *pullPoliciesHelpPtr {
			pullPoliciesCommand.PrintDefaults()
			os.Exit(0)
		} else if *pullPoliciesProduct == "" || *pullPoliciesUsernamePtr == "" || *pullPoliciesAuthoritiesPtr == "" {
			pullPoliciesCommand.PrintDefaults()
			os.Exit(1)
		}
		request := modules.PolicyPullRequest{*pullPoliciesIdPtr, *pullPoliciesExternalCodePtr, *pullPoliciesAgremmentPtr}
		auth := modules.Authorization{*pullPoliciesUsernamePtr, *pullPoliciesAuthoritiesPtr}
		modules.PullPolicies(*pullPoliciesProduct, request, auth, *pullPoliciesVerbosePtr)
	}

	if checkIbanCommand.Parsed() {
		if *checkIbanHelpPtr {
			checkIbanCommand.PrintDefaults()
			os.Exit(0)
		} else if *checkIbanValuePtr == "" {
			checkIbanCommand.PrintDefaults()
			os.Exit(1)
		}
		modules.CheckIban(*checkIbanCountryPtr, *checkIbanValuePtr, *checkIbanVerbosePtr)
	}
}

func usage() {
	fmt.Println(`
Usage: hodei-cli COMMAND [OPTIONS]")

Commands:
  read-person
  pull-countries
  pull-products
  pull-policies
  check-iban
  version
`)
}

