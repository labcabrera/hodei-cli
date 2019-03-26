package modules

import (
	"log"

	"github.com/labcabrera/hodei-cli/client"
)

func PullProfessions(verbose bool) {
	if verbose {
		log.Printf("Pulling professions from referential API")
	}
	client.SendMessage("cnp.referential", "profession.pull", "{}", verbose)
}
