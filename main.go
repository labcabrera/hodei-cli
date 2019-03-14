package main

import (
	"fmt"
	"log"
	"flag"
	"os"
	"github.com/streadway/amqp"
)

const version = "0.1.0"
const exchangeSepa = "cnp.sepa"
const exchangeReferential = "cnp.referential"

func main() {

	// Arguments
	pullCountries := flag.Bool("pull-countries", false, "Pull countries over referential API")
	pullProducts := flag.Bool("pull-products", false, "Pull products and agreements over referential API")
	pullPolicies := flag.Bool("pull-ppi-policies", false, "Pull PPI policies")
	pullPerson := flag.String("pull-person", "", "Pull person from referential API")
	pullLegal := flag.String("pull-legal", "", "Pull legal entity from referential API")
	policyId := flag.String("policy-id", "", "MongoDB Policy identifier")
	iban := flag.String("iban", "", "IBAN validation")
	verbose := flag.Bool("v", false, "Verbose")
	printVersion := flag.Bool("version", false, "Print version")
	flag.Parse()
	
	if(*printVersion) {
		fmt.Println("Hodei cli ", version)
		return
	}
	if(*iban != "") {
		sendIbanValdidationMessage(*iban, *verbose)
		return
	}
	if(*pullCountries) {
		sendPullCountryMessage(*verbose)
	}
	if(*pullProducts) {
		sendPullProductsMessage(*verbose)
	}
	if(*pullPerson != "") {
		//TODO
	}
	if(*pullLegal != "") {
		//TODO
	}
	if(*pullPolicies) {
		sendPullPoliciesMessage(*policyId, *verbose)
	}
}

func sendIbanValdidationMessage(iban string, verbose bool) {
	if(verbose) {
		log.Printf("Validating IBAN %s", iban)
	}
	body := `{"countryCode": "ESP","iban": "` + iban + `"}`
	sendMessage(exchangeSepa, "iban.validation", body, verbose)
}

func sendPullCountryMessage(verbose bool) {
	if(verbose) {
		log.Printf("Pulling countries from referential API")
	}
	sendMessage(exchangeReferential, "country.pull", "", verbose)
}

func sendPullProductsMessage(verbose bool) {
	if(verbose) {
		log.Printf("Pulling products and agreements from referential API")
	}
	sendMessage(exchangeReferential, "product.pull", "", verbose)
}

func sendPullPoliciesMessage(policyId string, verbose bool) {
	//TODO
}

func sendMessage(exchange string, routingKey string, body string, verbose bool) (err error) {
	amqpUri := "amqp://" + os.Getenv("APP_AMQP_URI")
	conn, err := amqp.Dial(amqpUri)
	if(err != nil) {
		log.Fatalf("%s: %s", "Error opening connection", err)
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if(err != nil) {
		log.Fatalf("%s: %s", "Error opening channel", err)
		return err
	}
	defer ch.Close()
	err = ch.Publish(
		exchange,
		routingKey,
		false,			// mandatory
		false,			// inmediate
		amqp.Publishing {
			ContentType:	"text/plain",
			Body:			[]byte(body),
		})
	if(err != nil) {
		log.Fatalf("%s: %s", "Error opening connection", err)
	} else if(verbose) {
		log.Printf("Sent message: %s", body)
	}
	return err
}
