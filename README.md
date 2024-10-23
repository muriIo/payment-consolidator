# Payment Consolidator

A project that use RabbitMQ to enqueue payment records read from a CSV file.

## How it works

We have a CSV file with some records mocking a payment list. Payments have types, which the _PRODUCER_ project reads and send them to the specific _RABBITMQ_ queue (pix queue, bank slip queue or card queue). Then, we have three different consumers, each consumer reads from its specific queue.

The exchange is direct.

## Instalation

This project uses _Docker_ to build a _RabbitMQ Management_ image and _Golang_, so:

### Prerequisites

Before you begin, ensure you have the following software installed:

- **Docker (with docker compose enabled)**
  - [Windows download here](https://docs.docker.com/desktop/install/windows-install/)
  - [Linux download here](https://docs.docker.com/desktop/install/linux/)
  - [Mac download here](https://docs.docker.com/desktop/install/mac-install/)
- **Go** (version 1.22.4 or higher) - [Download here](https://go.dev/doc/install)

### Step 1: Run the docker compose command

Open your terminal (and go to the root folder of the project) and run the following command to activate the RabbitMQ Management container:

```bash
docker-compose up -d
```

### Step 2: Run the code

We need to run the _consumers_ project in one terminal (because it is going to be consuming the queue until we cancel the process) and in another we run the _producer_ project to generate the data and send it to the queue.

#### First terminal

Go to the /consumers folder and run any of the folllowing commands:

- For the pix queue:

  ```bash
  go run main.go pix
  ```

- For the bank slip queue:

  ```bash
  go run main.go bank_slip
  ```

- For the card queue:

  ```bash
  go run main.go card
  ```

#### Second terminal

Go to the /producer folder and run the command:

```bash
go run producer.go MOCK_DATA.csv
```

### Step 3: See the result

You will see, in the format of _JSON_, all the records that the consumer project fetched from the specific queue that you provided.

## Nice things to notice

I went a little further adding the support of reading arguments from the command that runs the project. I thought that it would help if you want to try another CSV file (in the same format of columns, of course). Just provide other file in the command.

## Not stopping here

In the future, I want to introduce other features:

- Monitoring and logging
  - Enhance the perfomance
- Fault tolerance and retry mechanism
- Load balancing for consumers
- Parallel consumers per payment type
- Manual message acknowledgement and durability
- Batch processing
- Timeout and circuit breaker patterns
- Security enhancements (since we are dealing with financial data)
- Streamlining customer API integration

But that's a talk for another time.
