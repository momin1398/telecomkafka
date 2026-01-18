package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/momin1398/telecomkafka/internal/api"
	"github.com/momin1398/telecomkafka/internal/kafka"
	"github.com/momin1398/telecomkafka/internal/service"
)

// ---- Dummy producer (used when Kafka is unavailable)
type DummyProducer struct{}

func (d *DummyProducer) Publish(
	ctx context.Context,
	topic string,
	key, value []byte,
) error {
	log.Printf(
		"Dummy publish -> topic=%s key=%s value=%s\n",
		topic,
		string(key),
		string(value),
	)
	return nil
}

func main() {
	log.Println("Starting Collector API...")

	// ---------- Config ----------
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "redpanda:9092"
	}

	topic := os.Getenv("TELEMETRY_TOPIC")
	if topic == "" {
		topic = "metrics-topic"
	}

	// ---------- Producer selection ----------
	var producer kafka.MessageProducer

	kafkaProducer := kafka.NewKafkaProducer([]string{broker})

	// Health check: try producing a dummy message
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := kafkaProducer.Publish(
		ctx,
		topic,
		[]byte("health-check"),
		[]byte(`{"status":"ok"}`),
	)

	if err != nil {
		log.Printf(
			"Kafka unavailable (%v), using DummyProducer\n",
			err,
		)
		producer = &DummyProducer{}
	} else {
		log.Println("Kafka connected successfully")
		producer = kafkaProducer
	}

	// ---------- Service layer ----------
	telemetryService := service.NewTelemetryService(producer, topic)

	// ---------- HTTP handler ----------
	handler := api.NewHandler(telemetryService)

	mux := http.NewServeMux()
	mux.HandleFunc("/collect", handler.Collect)

	// ---------- Server ----------
	log.Println("Collector API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
