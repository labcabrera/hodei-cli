package modules

import (
	"log"
	"flag"
	"os"
	"github.com/streadway/amqp"
	"github.com/labcabrera/hodei-cli/client"
)

func CustomerSearch(id string, username string, authorities string, verbose bool) (res string, err error) {
	if verbose {
		log.Printf("Searching customer %s (%s:%s)", id, username, authorities)
	}
	if id == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}


	headers := amqp.Table{
		"App-Username"   : username,
		"App-Authorities": authorities,
	}
	body := `{"1":{"type":"person","reference":"` + id + `"}}`

	if(verbose) {
		log.Printf("Body: %s", body)
	}

	res, err = client.SendAndReceive("cnp.customer", "customer.search", body, headers, verbose)
	if(err != nil) {
		log.Fatalf("%s: %s", "Error reading person", err)
		return
	}
	
	return
}

