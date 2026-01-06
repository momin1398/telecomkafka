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
		GroupID: "metrics-service",
	})

	for {
		msg, _ := reader.ReadMessage(context.Background())
		log.Printf("Metrics stored: %s", string(msg.Value))
	}
}
