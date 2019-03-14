package main

import (
	"fmt"
	"log"
	"flag"
	"github.com/streadway/amqp"
)

const exchangeSepa = "cnp.sepa"
const routingKeyIbanValidation = "iban.validation"

func main() {
	fmt.Println("Hodei cli 0.1.0-SNAPSHOT")

	iban := flag.String("iban", "xxx", "IBAN validation")
	flag.Parse()

	fmt.Println("IBAN: ", *iban)
	
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	
	body := `{
		"countryCode": "ESP",
		"iban": "` + *iban + `"
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

	fmt.Printf("Ciao\n")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
