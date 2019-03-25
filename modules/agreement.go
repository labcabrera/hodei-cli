package modules

import (
	"log"

	"github.com/labcabrera/hodei-cli/client"
)

func PullAgreements(productId string, verbose bool) {
	if verbose {
		log.Printf("Pulling agreements from referential API")
	}
	body := `{}`
	client.SendMessage("cnp.referential", "agreement.pull", body, verbose)
}
