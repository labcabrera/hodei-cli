package modules

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/labcabrera/hodei-cli/client"
)

func CustomerSearch(id string, auth Authorization, verbose bool) {
	if(verbose) {
		log.Printf("Searching customer %s (%s)", id, auth)
	}

	headers := amqp.Table{
		"App-Username"   : auth.Username,
		"App-Authorities": auth.Authorities,
	}

	id = "5c82818bd601800001c95776"

	body := `{"` + id + `":{"type":"person","reference":"` + id + `"}}`

	client.SendMessageWithHeaders("cnp.customer", "customer.search", body, headers, verbose)
}