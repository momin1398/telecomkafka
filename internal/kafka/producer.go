package kafka

import (
	"context"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

// MessageProducer defines the interface for publishing messages
type MessageProducer interface {
	Publish(ctx context.Context, topic string, key, value []byte) error
}

// -----------------
// Real Kafka Producer
// -----------------

type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a KafkaProducer connected to the given brokers
func NewKafkaProducer(brokers []string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// Publish sends a message to the specified topic with key and value
func (p *KafkaProducer) Publish(ctx context.Context, topic string, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

// -----------------
// Dummy Producer (for testing without Kafka)
// -----------------

type DummyProducer struct{}

// NewDummyProducer creates a DummyProducer
func NewDummyProducer() *DummyProducer {
	return &DummyProducer{}
}

// Publish just prints the message instead of sending to Kafka
func (d *DummyProducer) Publish(ctx context.Context, topic string, key, value []byte) error {
	fmt.Printf("Dummy publish -> topic: %s key: %s value: %s\n", topic, key, string(value))
	return nil
}
