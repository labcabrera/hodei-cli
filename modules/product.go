package modules

import (
	"log"

	"github.com/labcabrera/hodei-cli/client"
)

func PullProducts(verbose bool) {
	if verbose {
		log.Printf("Pulling products from referential API")
	}
	client.SendMessage("cnp.referential", "product.pull", "", verbose)
}
