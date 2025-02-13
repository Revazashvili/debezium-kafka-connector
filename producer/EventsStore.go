package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/revazashvili/debezium-kafka-connector/events"
)

const sql = `INSERT INTO outbox.outbox_messages (payload, timestamp, aggregate_id, type) VALUES ($1, $2, $3, $4)`

type EventsStore struct {
	p *pgxpool.Pool
}

func NewEventsStore(p *pgxpool.Pool) EventsStore {
	return EventsStore{p: p}
}

func (es EventsStore) Save(messages []events.OutboxMessage, ctx context.Context) error {

	conn, err := es.p.Acquire(ctx)
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
