package iban

import (
	"log"
	"github.com/labcabrera/hodei-cli/client"
)

const exchangeSepa = "cnp.sepa"

func ProcessIban(iban string, verbose bool) {
	if(verbose) {
		log.Printf("Validating IBAN %s", iban)
	}
	body := `{"countryCode": "ESP","iban": "` + iban + `"}`
	client.SendMessage(exchangeSepa, "iban.validation", body, verbose)
}