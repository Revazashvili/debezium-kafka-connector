package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
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

	createSchema = "CREATE SCHEMA IF NOT EXISTS outbox"
)

type DbInitializer struct {
	p *pgxpool.Pool
}

func NewDbInitializer(p *pgxpool.Pool) *DbInitializer {
	return &DbInitializer{p: p}
}

func (di DbInitializer) Init(ctx context.Context) error {
	c, err := di.p.Acquire(ctx)
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
