package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/labcabrera/hodei-cli/modules"
)

const version = "0.3.0"
const versionCmd = "version"

func main() {

	// Check command argument
	if len(os.Args) < 2 {
		usage()
		return
	}

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

	checkIbanOptions := modules.CheckIbanOptions{}
	checkIbanFlagSet := modules.CheckIbanFlagSet(&checkIbanOptions)

	rand.Seed(time.Now().UTC().UnixNano())

	switch os.Args[1] {
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
	case modules.CheckIbanCmd:
		checkIbanFlagSet.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
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

	if checkIbanFlagSet.Parsed() {
		if checkIbanOptions.Help {
			checkIbanFlagSet.PrintDefaults()
			os.Exit(0)
		}
		modules.CheckIban(&checkIbanOptions)
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
  ` + modules.PullPoliciesCmd + `  
  ` + modules.PullProfessionsCmd + `
  ` + modules.CheckIbanCmd + `
  ` + versionCmd)
}
