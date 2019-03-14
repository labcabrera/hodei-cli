package main

import (
	"fmt"
	"log"
	"flag"
	"github.com/streadway/amqp"
)

const version = "0.1.0-SNAPSHOT"
const exchangeSepa = "cnp.sepa"
const routingKeyIbanValidation = "iban.validation"

func main() {
	// Argument variables
	pullCountries := flag.Bool("pull-countries", false, "Pull countries over referential API")
	pullProducts := flag.Bool("pull-products", false, "Pull products and agreements over referential API")
	pullPolicies := flag.Bool("pull-ppi-policies", false, "Pull PPI policies")
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
	if(*pullPolicies) {
		sendPullPoliciesMessage(*policyId, *verbose)
	}
}

func sendIbanValdidationMessage(iban string, verbose bool) {
	if(verbose) {
		fmt.Println("Validating IBAN", iban)
	}
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	
	body := `{
		"countryCode": "ESP",
		"iban": "` + iban + `"
	}`
	err = ch.Publish(
		exchangeSepa,						// exchange
		routingKeyIbanValidation,			// routing key
		false,								// mandatory
		false,								// immediate
		amqp.Publishing {
			ContentType:	"text/plain",
			Body:			[]byte(body),
		})
	
	failOnError(err, "Failed to publish a message")
}

func sendPullCountryMessage(verbose bool) {
	if(verbose) {
		fmt.Println("Pulling countries from referential API")	
	}
	fmt.Println("Not implemented: pull countries")
}

func sendPullProductsMessage(verbose bool) {
}

func sendPullPoliciesMessage(policyId string, verbose bool) {
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
