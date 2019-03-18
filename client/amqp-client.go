package client

import(
	"log"
	"os"
	"github.com/streadway/amqp"
)

func SendMessage(exchange string, routingKey string, body string, verbose bool) (err error) {
	return SendMessageWithHeaders(exchange, routingKey, body, nil, verbose)
}

func SendMessageWithHeaders(exchange string, routingKey string, body string, headers amqp.Table, verbose bool) (err error) {
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