package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
)

const (
	broker = "localhost:29092"
)

var topics = []string{"products.ProductCreatedEvent", "products.ProductNameUpdatedEvent", "products.ProductDeletedEvent"}

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "products",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return
	}
	defer c.Close()

	err = c.SubscribeTopics(topics, nil)

	if err != nil {
		fmt.Printf("Failed to subscribe to topics: %s\n", err)
		return
	}

	fmt.Printf("Subscribed to topics: %s\n", topics)
	mc := 0

	osc := make(chan os.Signal, 1)
	signal.Notify(osc, os.Interrupt)
	for {
		select {
		case <-osc:
			fmt.Printf("done")
			return
		default:
			message, err := c.ReadMessage(-1)

			if err != nil {
				fmt.Printf("Error reading message: %s\n", err)
				continue
			}

			mc++
			fmt.Println(string(message.Value))
			fmt.Printf("Consumed Message count: %v \n", mc)
		}
	}
}
