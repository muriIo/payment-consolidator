package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Payment struct {
	ID              int64  `json:"id"`
	Customer_id     string `json:"customer_id"`
	Costumer_name   string `json:"customer_name"`
	Description     string `json:"description"`
	Amount          string `json:"amount"`
	Type_of_payment string `json:"type_of_payment"`
}

type Payments []Payment

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func readCSV(name string) [][]string {
	file, err := os.Open(name)

	failOnError(err, "Error while reading the csv file. Check if the file exists and try again.")

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	failOnError(err, "Error reading the csv's content.")

	return records
}

func normalizeRecords(records [][]string) Payments {
	var payments Payments

	for i, record := range records {
		if i == 0 {
			continue
		}

		id, err := strconv.ParseInt(record[0], 10, 64)

		failOnError(err, "Error converting id of the record.")

		payment := Payment{
			ID:              id,
			Customer_id:     record[1],
			Costumer_name:   record[2],
			Description:     record[3],
			Amount:          record[4],
			Type_of_payment: record[5],
		}

		payments = append(payments, payment)
	}

	return payments
}

func (payments Payments) toJson() []byte {
	message, err := json.MarshalIndent(payments, "", "  ")

	failOnError(err, "Error converting struct to JSON")

	return message
}

func (payments Payments) filterByTypeOfPayment(typeOfPayment string) Payments {
	var filtered Payments

	for _, payment := range payments {
		if payment.Type_of_payment == typeOfPayment {
			filtered = append(filtered, payment)
		}
	}

	return filtered
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No correct arguments provided. One is expected.")
	}

	records := readCSV(os.Args[1])
	payments := normalizeRecords(records)
	pix_payment := payments.filterByTypeOfPayment("pix")
	bank_slip_payment := payments.filterByTypeOfPayment("bank_slip")
	card_payment := payments.filterByTypeOfPayment("card")
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

	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = ch.PublishWithContext(ctx,
		exchange_name,
		"pix",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        pix_payment.toJson(),
		})

	failOnError(err, "Failed to publish the pix records")

	err = ch.PublishWithContext(ctx,
		exchange_name,
		"bank_slip",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bank_slip_payment.toJson(),
		})

	failOnError(err, "Failed to publish the bank slip records")

	err = ch.PublishWithContext(ctx,
		exchange_name,
		"card",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        card_payment.toJson(),
		})

	failOnError(err, "Failed to publish the card records")
}
