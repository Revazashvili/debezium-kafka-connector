package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/revazashvili/debezium-kafka-connector/events"
)

const (
	createTable = `
CREATE TABLE IF NOT EXISTS outbox.outbox_messages (
    id SERIAL PRIMARY KEY,
    payload jsonb NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    aggregate_id BIGINT NOT NULL,
    type VARCHAR(255) NOT NULL
);`

	createSchema = `CREATE SCHEMA IF NOT EXISTS outbox`

	sql = `INSERT INTO outbox.outbox_messages (payload, timestamp, aggregate_id, type) VALUES ($1, $2, $3, $4)`
)

type postgresOutboxStore struct {
	p *pgxpool.Pool
}

func NewPostgresOutboxStore(p *pgxpool.Pool) OutboxStore {
	return &postgresOutboxStore{p: p}
}

func (pos *postgresOutboxStore) Init(ctx context.Context) error {
	c, err := pos.p.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	tx, err := c.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, createSchema)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, createTable)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (pos *postgresOutboxStore) Save(ctx context.Context, messages []events.OutboxMessage) error {
	conn, err := pos.p.Acquire(ctx)
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
