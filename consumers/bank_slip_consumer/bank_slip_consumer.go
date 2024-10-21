package bank_slip_consumer

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
	fmt.Printf("Initialzing bank slip consumer...\n")

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

	bank_slip_q, err := ch.QueueDeclare(
		"bank_slip",
		false,
		false,
		true,
		false,
		nil,
	)

	failOnError(err, "Failed to create the pix queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s", bank_slip_q.Name, exchange_name, "bank_slip")

	err = ch.QueueBind(
		bank_slip_q.Name,
		"bank_slip",
		exchange_name,
		false,
		nil,
	)

	failOnError(err, "Failed to bind the bank slip queue")

	bank_slip_payments, err := ch.Consume(
		bank_slip_q.Name,
		"bank_slip_payments_consumer",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer to the bank slip queue")

	var forever chan struct{}

	go func() {
		for d := range bank_slip_payments {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for payments by bank slip.")
	<-forever
}
