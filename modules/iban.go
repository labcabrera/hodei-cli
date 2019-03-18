package modules

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/labcabrera/hodei-cli/client"
)

const exchangeSepa = "cnp.sepa"

func CheckIban(country string, iban string, verbose bool) (res string, err error) {
	if(verbose) {
		log.Printf("Validating IBAN %s", iban)
	}
	headers := amqp.Table{
		"App-Source":		"hodei-cli",
	}
	body := `{"countryCode": "` + country + `","iban": "` + iban + `"}`
	res, err = client.SendAndReceive("cnp.sepa", "iban.validation", body, headers, verbose)
	return
}