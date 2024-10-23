package card_consumer

import (
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Listen(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Initialzing card consumer...\n")

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

	card_q, err := ch.QueueDeclare(
		"card",
		false,
		false,
		true,
		false,
		nil,
	)

	failOnError(err, "Failed to create the pix queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s", card_q.Name, exchange_name, "card")

	err = ch.QueueBind(
		card_q.Name,
		"card",
		exchange_name,
		false,
		nil,
	)

	failOnError(err, "Failed to bind the card queue")

	card_payments, err := ch.Consume(
		card_q.Name,
		"card_payments_consumer",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer to the card queue")

	var forever chan struct{}

	go func() {
		for d := range card_payments {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for payments by card.")
	<-forever
}
