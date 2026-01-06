package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "device-telemetry",
		GroupID: "audit-service",
	})

	for {
		msg, _ := reader.ReadMessage(context.Background())
		log.Printf("AUDIT: %s", string(msg.Value))
	}
}
