package events

import "time"

type OutboxMessage struct {
	Id          int
	Payload     string
	Timestamp   time.Time
	AggregateId int
	Type        string
}
