package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Printf("Hodei cli 0.1.0-SNAPSHOT\n")
	
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	
	body := `{
		"countryCode": "ESP",
		"iban": "invalid-iban"
	}`
	err = ch.Publish(
		"cnp.sepa",							// exchange
		"iban.validation",					// routing key
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
