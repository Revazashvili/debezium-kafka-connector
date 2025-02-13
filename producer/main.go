package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/revazashvili/debezium-kafka-connector/events"
	"os"
)

const dbURL = "postgres://user:pass@localhost:5432/debezium"

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	err = createSchemaIfNotExists(pool, ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to create schema: %v\n", err)
	}

	err = createTableIfNotExists(pool, ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
	}

	es := events.GenerateEvents()
	err = insertOutboxMessages(pool, es, ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to insert outbox messages: %v\n", err)
		return
	}

	fmt.Println("Successfully inserted outbox messages")

	//for {
	//	time.Sleep(time.Second * 3)
	//
	//	es := events.GenerateEvents()
	//	err = insertOutboxMessages(pool, es, ctx)
	//	if err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "Unable to insert outbox messages: %v\n", err)
	//		return
	//	}
	//
	//	fmt.Println("Successfully inserted outbox messages")
	//}
}

func insertOutboxMessages(pool *pgxpool.Pool, messages []events.OutboxMessage, ctx context.Context) error {
	sql := `INSERT INTO outbox.outbox_messages (payload, timestamp, aggregate_id, type) VALUES ($1, $2, $3, $4)`

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)

	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, msg := range messages {
		_, err := tx.Exec(ctx, sql, msg.Payload, msg.Timestamp, msg.AggregateId, msg.Type)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func createTableIfNotExists(pool *pgxpool.Pool, ctx context.Context) error {
	createTable := `
CREATE TABLE IF NOT EXISTS outbox.outbox_messages (
    id SERIAL PRIMARY KEY,
    payload jsonb NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    aggregate_id BIGINT NOT NULL,
    type VARCHAR(255) NOT NULL
);`

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, createTable)
	if err != nil {
		return err
	}

	return nil
}

func createSchemaIfNotExists(pool *pgxpool.Pool, ctx context.Context) error {
	createSchema := "CREATE SCHEMA IF NOT EXISTS outbox"

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, createSchema)
	if err != nil {
		return err
	}

	return nil
}
