package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	// ---------- Config ----------
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "redpanda:9092"
	}

	topic := os.Getenv("TOPIC")
	if topic == "" {
		log.Fatal("TOPIC env variable not set")
	}

	log.Printf("starting consumer for topic: %s\n", topic)

	// ---------- Kafka Reader ----------
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  topic + "-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	defer reader.Close()

	// ---------- Context with graceful shutdown ----------
	ctx, cancel := context.WithCancel(context.Background())

	// Handle SIGTERM / CTRL+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		log.Println("shutdown signal received, closing consumer...")
		cancel()
	}()

	// ---------- Consume loop ----------
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Println("consumer stopped:", err)
			return
		}

		log.Printf(
			"✅ topic=%s partition=%d offset=%d key=%s value=%s\n",
			msg.Topic,
			msg.Partition,
			msg.Offset,
			string(msg.Key),
			string(msg.Value),
		)
	}
}
