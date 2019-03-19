package client

import(
	"log"
	"os"
	"time"
	"math/rand"
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

	msg := amqp.Publishing{
		DeliveryMode:	amqp.Persistent,
		Timestamp:		time.Now(),
		ContentType:	"text/plain",
		Headers:		headers,
		Body:			[]byte(body),
	}

	err = ch.Publish(
		exchange,
		routingKey,
		false,			// mandatory
		false,			// inmediate
		msg)
	if(err != nil) {
		log.Fatalf("%s: %s", "Error opening connection", err)
	} else if(verbose) {
		log.Printf("Sent message")
	}
	return err
}

func SendAndReceive(exchange string, routingKey string, body string, headers amqp.Table, verbose bool) (res string, err error) {
	amqpUri := "amqp://" + os.Getenv("APP_AMQP_URI")
	conn, err := amqp.Dial(amqpUri)
	if(err != nil) {
		log.Fatalf("%s: %s", "Error opening connection", err)
		return "", err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if(err != nil) {
		log.Fatalf("%s: %s", "Error opening channel", err)
		return "", err
	}
	defer ch.Close()
	if(verbose) {
		log.Printf("Sending message: %s", body)
	}

	q, err := ch.QueueDeclare(
		"",						// name
		false,					// durable
		false,					// delete when usused
		true,					// exclusive
		false,					// noWait
		nil,					// arguments
		)

	msgs, err := ch.Consume(
		q.Name,					// queue
		"",						// consumer
		true,					// auto-ack
		false,					// exclusive
		false,					// no-local
		false,					// no-wait
		nil,					// args
		)

	corrId := randomString(32)

	if(verbose) {
		log.Printf("Using correlation id: %s", corrId)
	}

	err = ch.Publish(
		exchange,				// exchange
		routingKey,				// routing key
		false,					// mandatory
		false,					// immediate
		amqp.Publishing{
			ContentType:	"text/plain",
			CorrelationId:	corrId,
			ReplyTo:		q.Name,
			Headers:		headers,
			Body:			[]byte(body),
		})

	for d := range msgs {
		if corrId == d.CorrelationId {
			res = string(d.Body)
			break
		} else {
			log.Printf("Test %s", d)
		}
	}
	return
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
			bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
