package pix_consumer

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Listen() {
	fmt.Printf("Initialzing pix consumer...\n")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()

	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	exchange_name := "payments_direct"

	err = ch.ExchangeDeclare(
		exchange_name,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	pix_q, err := ch.QueueDeclare(
		"pix",
		false,
		false,
		true,
		false,
		nil,
	)

	failOnError(err, "Failed to create the pix queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s", pix_q.Name, exchange_name, "pix")

	err = ch.QueueBind(
		pix_q.Name,
		"pix",
		exchange_name,
		false,
		nil,
	)

	failOnError(err, "Failed to bind the pix queue")

	pix_payments, err := ch.Consume(
		pix_q.Name,
		"pix_payments_consumer",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer to the pix queue")

	var forever chan struct{}

	go func() {
		for d := range pix_payments {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for payments by pix.")
	<-forever
}
