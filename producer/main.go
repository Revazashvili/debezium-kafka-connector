package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/revazashvili/debezium-kafka-connector/events"
	"os"
	"time"
)

const dbURL = "postgres://postgres:mysecretpassword@localhost:5433/debezium"

func main() {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	err = createSchemaIfNotExists(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create schema: %v\n", err)
	}

	err = createTableIfNotExists(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
	}

	for {
		time.Sleep(time.Second * 3)

		es := events.GenerateEvents()
		err = insertOutboxMessages(conn, es)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert outbox messages: %v\n", err)
		}

		fmt.Println("Successfully inserted outbox messages")
	}
}

func insertOutboxMessages(conn *pgx.Conn, messages []events.OutboxMessage) error {
	sql := `INSERT INTO outbox.outbox_messages (payload, timestamp, aggregate_id, type) VALUES ($1, $2, $3, $4)`

	batch := &pgx.Batch{}
	for _, msg := range messages {
		batch.Queue(sql, msg.Payload, msg.Timestamp, msg.AggregateId, msg.Type)
	}

	br := conn.SendBatch(context.Background(), batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to close batch: %v\n", err)
		}
	}(br)

	// Consume batch results
	for range messages {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func createTableIfNotExists(conn *pgx.Conn) error {
	createTable := `
CREATE TABLE IF NOT EXISTS outbox.outbox_messages (
    id SERIAL PRIMARY KEY,
    payload TEXT NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    aggregate_id BIGINT NOT NULL,
    type VARCHAR(255) NOT NULL
);`
	_, err := conn.Exec(context.Background(), createTable)
	if err != nil {
		return err
	}

	return nil
}

func createSchemaIfNotExists(con *pgx.Conn) error {
	createSchema := "CREATE SCHEMA IF NOT EXISTS outbox"

	_, err := con.Exec(context.Background(), createSchema)
	if err != nil {
		return err
	}

	return nil
}
