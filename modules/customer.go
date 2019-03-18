package modules

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/labcabrera/hodei-cli/client"
)

func CustomerSearch(id string, customerType string, auth Authorization, verbose bool) {
	if(verbose) {
		log.Printf("Searching customer %s (%s)", id, auth)
	}
	headers := amqp.Table{
		"App-Username"   : auth.Username,
		"App-Authorities": auth.Authorities,
	}
	body := `{"` + id + `":{"type":"` + customerType + `","reference":"` + id + `"}}`
	client.SendMessageWithHeaders("cnp.customer", "customer.search", body, headers, verbose)
}