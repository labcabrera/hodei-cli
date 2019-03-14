package main

import (
	"fmt"
	"log"
	"flag"
	"os"
	"encoding/json"
	"github.com/streadway/amqp"
)

const version = "0.2.0-SNAPSHOT"
const exchangeSepa = "cnp.sepa"
const exchangeReferential = "cnp.referential"

type PolicyPullRequest struct {
	entityId			string
	externalCode		string
	agreementId			string
}

type Authorization struct {
	username			string
	authorities			string
}

func main() {

	// Arguments
	pullCountries  := flag.Bool("pull-countries",  false, "Pull countries over referential API")
	pullProducts   := flag.Bool("pull-products",   false, "Pull products and agreements over referential API")
	pullPolicies   := flag.Bool("pull-policies",   false, "Pull policies")
	pullPerson     := flag.Bool("pull-person",     false, "Pull person from referential API")
	pullLegal      := flag.Bool("pull-legal",      false, "Pull legal entity from referential API")

	id             := flag.String("id",            "",    "Entity MongoDB identifier")
	externalCode   := flag.String("external-code", "",    "Entity external code")
	product        := flag.String("product",       "",    "Product name")
	agreement      := flag.String("agreement",     "",    "Agreement external code")
	iban           := flag.String("iban",          "",    "IBAN validation")

	username       := flag.String("u",             "",    "Username")
	authorities    := flag.String("a",             "",    "Authorities (coma separated list)")
	verbose        := flag.Bool("v",               false, "Verbose")
	printVersion   := flag.Bool("version",         false, "Print version")
	
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
	if(*pullPerson) {
		//TODO
		fmt.Println("Not implemented")
	}
	if(*pullLegal) {
		//TODO
		fmt.Println("Not implemented")
	}
	if(*pullPolicies) {
		request := PolicyPullRequest{*id, *externalCode, *agreement}
		auth := Authorization{*username, *authorities}
		sendPullPoliciesMessage(*product, request, auth, *verbose)
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

func sendPullPoliciesMessage(product string, request PolicyPullRequest, auth Authorization, verbose bool) {
	switch(product) {
	case "":
		fmt.Println("Required argument product")
		return
	case "ppi":
		log.Printf("Agreement: %s", request.agreementId)

		bodyBinary, err := json.Marshal(request)
		if(err != nil) {
			log.Fatalf("%s: %s", "Error marshalling request", err)
			return
		}
		body := string(bodyBinary)

		//TODO
		body = `{"agreementId":"` + request.agreementId + `"}`		
		headers := amqp.Table{
			"App-Username"   : auth.username,
			"App-Authorities": auth.authorities,
		}
		sendMessageWithHeaders("ppi.referential", "policy.pull", body, headers, verbose)
	default:
		log.Fatalf("Unknown product %s", product)
	}
}

func sendMessage(exchange string, routingKey string, body string, verbose bool) (err error) {
	return sendMessageWithHeaders(exchange, routingKey, body, nil, verbose)
}

func sendMessageWithHeaders(exchange string, routingKey string, body string, headers amqp.Table, verbose bool) (err error) {
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
	if(verbose) {
		log.Printf("Sending message: %s", body)
	}
	err = ch.Publish(
		exchange,
		routingKey,
		false,			// mandatory
		false,			// inmediate
		amqp.Publishing {
			ContentType:	"text/plain",
			Body:			[]byte(body),
			Headers:		headers,
		})
	if(err != nil) {
		log.Fatalf("%s: %s", "Error opening connection", err)
	} else if(verbose) {
		log.Printf("Sent message")
	}
	return err
}
