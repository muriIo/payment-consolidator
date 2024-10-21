module payment-consolidator/consumer

go 1.22.4

replace payment-consolidator/consumers/pix_consumer => ./pix_consumer

replace payment-consolidator/consumers/card_consumer => ./card_consumer

replace payment-consolidator/consumers/bank_slip_consumer => ./bank_slip_consumer

require (
	payment-consolidator/consumers/bank_slip_consumer v0.0.0-00010101000000-000000000000
	payment-consolidator/consumers/card_consumer v0.0.0-00010101000000-000000000000
	payment-consolidator/consumers/pix_consumer v0.0.0-00010101000000-000000000000
)

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect
