package modules

import (
	"log"
	"github.com/labcabrera/hodei-cli/client"
)

func PullCountries(verbose bool) {
	if(verbose) {
		log.Printf("Pulling countries from referential API")
	}
	client.SendMessage("cnp.referential", "country.pull", "", verbose)
}
