# Outbox Event Publishing with [Debezium](https://debezium.io/) and [Kafka](https://kafka.apache.org/)

This project demonstrates publishing outbox events using [Debezium](https://debezium.io/) to a Kafka topic. It consists of:

- **Producer Service**: Writes outbox events to PostgreSQL every 3 seconds.
- **Debezium Initializer**: Configures Debezium via an HTTP `PUT` request.
- **Consumer Service**: Listens for messages from Kafka.

## Running the Project

### 1. Start Infrastructure Services

Run the following command to start PostgreSQL, Kafka, and Debezium:
```sh
docker-compose up -d
```

### 2. Run Golang Services

Each service should be started separately:
```sh
# Start Producer Service
cd /producer
go run main.go

# Start Debezium Initializer
cd /debezium-initializer
go run main.go

# Start Consumer Service
cd /consumer
go run main.go
```

### 3. Verify the Setup

- Ensure Kafka, PostgreSQL, and Debezium are running by checking the logs:
  ```sh
  docker-compose logs -f
  ```
- Check if the producer inserts outbox events into PostgreSQL, producers logs message when successfully inserts messages.
- Confirm the consumer receives and processes events from Kafka, check console log of running service and compare published and consumed message count.

## Cleanup

To stop and remove the infrastructure services:
```sh
docker-compose down
```
