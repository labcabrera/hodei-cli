package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/labcabrera/hodei-cli/modules"
)

const version = "0.7.0-SNAPSHOT"
const versionCmd = "version"

func main() {

	if len(os.Args) < 2 {
		usage()
		return
	}

	cmd := os.Args[1]

	var moduleMap map[string]modules.HodeiCliModule
	moduleMap = make(map[string]modules.HodeiCliModule)

	moduleMap[modules.ListScheduledActionsCmd] = modules.ListScheduledActionsModule{}
	moduleMap[modules.MongoResetCmd] = modules.MongoResetModule{}

	module, check := moduleMap[cmd]

	if !check {
		usage()
		os.Exit(1)
	} else {
		module.Execute(os.Args[2:])
		os.Exit(0)
	}

	//TODO convertir a modulos

	// Parse cmd options
	customerSearchOptions := modules.CustomerSearchOptions{}
	customerSearchFlagSet := modules.CustomerSearchFlagSet(&customerSearchOptions)

	pullCountriesOptions := modules.PullCountriesOptions{}
	pullCountriesFlagSet := modules.PullCountriesFlagSet(&pullCountriesOptions)

	pullProductsOptions := modules.PullProductsOptions{}
	pullProductsFlagSet := modules.PullProductsFlagSet(&pullProductsOptions)

	pullAgreementsOptions := modules.PullAgreementsOptions{}
	pullAgreementsFlagSet := modules.PullAgreementsFlagSet(&pullAgreementsOptions)

	pullNetworksOptions := modules.PullNetworksOptions{}
	pullNetworksFlagSet := modules.PullNetworksFlagSet(&pullNetworksOptions)

	pullCustomersOptions := modules.PullCustomerOptions{}
	pullCustomersFlagSet := modules.PullCustomersFlagSet(&pullCustomersOptions)

	pullProfessionsOptions := modules.PullProfessionsOptions{}
	pullProfessionsFlagSet := modules.PullProfessionsFlagSet(&pullProfessionsOptions)

	pullPoliciesOptions := modules.PullPoliciesOptions{}
	pullPoliciesFlagSet := modules.PullPoliciesFlagSet(&pullPoliciesOptions)

	pullOrdersOptions := modules.PullOrdersOptions{}
	pullOrdersFlagSet := modules.PullOrdersFlagSet(&pullOrdersOptions)

	pullCoveragesOptions := modules.PullCoveragesOptions{}
	pullCoveragesFlagSet := modules.PullCoveragesFlagSet(&pullCoveragesOptions)

	pullClaimsOptions := modules.PullClaimsOptions{}
	pullClaimsFlagSet := modules.PullClaimsFlagSet(&pullClaimsOptions)

	checkIbanOptions := modules.CheckIbanOptions{}
	checkIbanFlagSet := modules.CheckIbanFlagSet(&checkIbanOptions)

	signatureRequestOptions := modules.SignatureRequestOptions{}
	signatureRequestFlagSet := modules.SignatureRequestFlagSet(&signatureRequestOptions)

	rand.Seed(time.Now().UTC().UnixNano())

	switch cmd {
	case versionCmd:
		fmt.Println("Hodei cli", version)
		os.Exit(0)
	case modules.CustomerSearchCmd:
		customerSearchFlagSet.Parse(os.Args[2:])
	case modules.PullCountriesCmd:
		pullCountriesFlagSet.Parse(os.Args[2:])
	case modules.PullProductsCmd:
		pullProductsFlagSet.Parse(os.Args[2:])
	case modules.PullAgreementsCmd:
		pullAgreementsFlagSet.Parse(os.Args[2:])
	case modules.PullCustomersCmd:
		pullCustomersFlagSet.Parse(os.Args[2:])
	case modules.PullNetworksCmd:
		pullNetworksFlagSet.Parse(os.Args[2:])
	case modules.PullProfessionsCmd:
		pullProfessionsFlagSet.Parse(os.Args[2:])
	case modules.PullPoliciesCmd:
		pullPoliciesFlagSet.Parse(os.Args[2:])
	case modules.PullOrdersCmd:
		pullOrdersFlagSet.Parse(os.Args[2:])
	case modules.PullCoveragesCmd:
		pullCoveragesFlagSet.Parse(os.Args[2:])
	case modules.PullClaimsCmd:
		pullClaimsFlagSet.Parse(os.Args[2:])
	case modules.CheckIbanCmd:
		checkIbanFlagSet.Parse(os.Args[2:])
	case modules.SignatureRequestCmd:
		signatureRequestFlagSet.Parse(os.Args[2:])
	default:
		fmt.Printf("%s: '%s' is not a hodei-cli command.\n", os.Args[0], cmd)
		usage()
		os.Exit(1)
	}

	if customerSearchFlagSet.Parsed() {
		if customerSearchOptions.Help {
			customerSearchFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.CustomerSearch(&customerSearchOptions)
	}

	if pullCountriesFlagSet.Parsed() {
		if pullCountriesOptions.Help {
			pullCountriesFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullCountries(&pullCountriesOptions)
	}

	if pullProductsFlagSet.Parsed() {
		if pullProductsOptions.Help {
			pullProductsFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullProducts(&pullProductsOptions)
	}

	if pullAgreementsFlagSet.Parsed() {
		if pullAgreementsOptions.Help {
			pullAgreementsFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullAgreements(&pullAgreementsOptions)
	}

	if pullNetworksFlagSet.Parsed() {
		if pullNetworksOptions.Help {
			pullNetworksFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullNetworks(&pullNetworksOptions)
	}

	if pullCustomersFlagSet.Parsed() {
		if pullCustomersOptions.Help {
			pullCustomersFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullCustomers(&pullCustomersOptions)
	}

	if pullProfessionsFlagSet.Parsed() {
		if pullProfessionsOptions.Help {
			pullProfessionsFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullProfessions(&pullProfessionsOptions)
	}

	if pullPoliciesFlagSet.Parsed() {
		if pullPoliciesOptions.Help {
			pullPoliciesFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullPolicies(&pullPoliciesOptions)
	}

	if pullOrdersFlagSet.Parsed() {
		if pullOrdersOptions.Help {
			pullOrdersFlagSet.PrintDefaults()
			os.Exit(0)
		}
		res, err := modules.PullOrders(&pullOrdersOptions)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(res)
		}
	}

	if pullCoveragesFlagSet.Parsed() {
		if pullCoveragesOptions.Help {
			pullCoveragesFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullCoverages(&pullCoveragesOptions)
	}

	if pullClaimsFlagSet.Parsed() {
		if pullClaimsOptions.Help {
			pullClaimsFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.PullClaims(&pullClaimsOptions)
	}

	if checkIbanFlagSet.Parsed() {
		if checkIbanOptions.Help {
			checkIbanFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.CheckIban(&checkIbanOptions)
	}

	if signatureRequestFlagSet.Parsed() {
		if signatureRequestOptions.Help {
			signatureRequestFlagSet.PrintDefaults()
			os.Exit(0)
		}
		res, err := modules.SignatureRequest(&signatureRequestOptions)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(res)
		}
	}
}

func usage() {
	fmt.Println(`
Usage: hodei-cli COMMAND [OPTIONS]")

Commands:
  ` + modules.CustomerSearchCmd + `
  ` + modules.PullCountriesCmd + `
  ` + modules.PullProductsCmd + `
  ` + modules.PullAgreementsCmd + `
  ` + modules.PullNetworksCmd + `
  ` + modules.PullCustomersCmd + `
  ` + modules.PullProfessionsCmd + `
  ` + modules.PullPoliciesCmd + `
  ` + modules.PullOrdersCmd + `
  ` + modules.PullCoveragesCmd + `
  ` + modules.PullClaimsCmd + `
  ` + modules.CheckIbanCmd + `
  ` + modules.MongoResetCmd + `
  ` + modules.SignatureRequestCmd + `
  ` + versionCmd)
}
