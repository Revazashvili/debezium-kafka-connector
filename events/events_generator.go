package events

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/jaswdr/faker"
)

func GenerateEvents() []OutboxMessage {
	f := faker.New()
	events := make([]OutboxMessage, 10)

	for i := 0; i < len(events); i++ {
		e := generateRandomEvent(f)
		b, err := json.Marshal(e)

		if err != nil {
			panic(err)
		}

		events[i] = OutboxMessage{
			Id:          f.Int(),
			Type:        e.GetType(),
			AggregateId: e.GetAggregateId(),
			Timestamp:   time.Now(),
			Payload:     string(b),
		}
	}
	return events
}

func generateRandomEvent(f faker.Faker) Event {
	rand.Seed(time.Now().UnixNano())
	eventType := rand.Intn(3)

	switch eventType {
	case 0:
		return ProductCreatedEvent{
			Id:       f.Int(),
			Name:     f.Person().FirstName(),
			ImageUrl: f.Internet().URL(),
			Price:    f.Int(),
		}
	case 1:
		return ProductNameUpdatedEvent{
			Id:   f.Int(),
			Name: f.Person().FirstName(),
		}
	case 2:
		return ProductDeletedEvent{
			Id: f.Int(),
		}
	default:
		panic("invalid event type")
	}
}
