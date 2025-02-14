package main

import (
	"context"
	"github.com/revazashvili/debezium-kafka-connector/events"
)

type OutboxStore interface {
	Init(ctx context.Context) error
	Save(ctx context.Context, messages []events.OutboxMessage) error
}
