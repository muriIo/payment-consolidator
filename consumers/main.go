package main

import (
	"fmt"
	"log"
	"os"

	"payment-consolidator/consumers/bank_slip_consumer"
	"payment-consolidator/consumers/card_consumer"
	"payment-consolidator/consumers/pix_consumer"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("Usage: %s [pix] or [bank_slip] or [card]", os.Args[0])
		os.Exit(0)
	}

	payment_mapper := map[string]func(){
		"pix":       pix_consumer.Listen,
		"bank_slip": bank_slip_consumer.Listen,
		"card":      card_consumer.Listen,
	}

	key := os.Args[1]

	if processFunc, ok := payment_mapper[key]; ok {
		processFunc()
	} else {
		fmt.Printf("Unknown payment method: %s\n", key)
	}
}
