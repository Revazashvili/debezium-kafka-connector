package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/revazashvili/debezium-kafka-connector/events"
	"os"
	"time"
)

const dbURL = "postgres://user:pass@localhost:5432/debezium"

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	pos := NewPostgresOutboxStore(pool)
	err = pos.Init(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to setup table: %v\n", err)
	}

	for {
		time.Sleep(time.Second * 3)

		es := events.GenerateEvents()
		err = pos.Save(ctx, es)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to insert outbox messages: %v\n", err)
			return
		}

		fmt.Println("Successfully inserted outbox messages")
	}
}
