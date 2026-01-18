package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	kafka1 "github.com/momin1398/telecomkafka/internal/kafka"
)

type MetricsEvent struct {
	CPU float64 `json:"cpu"`
	RAM float64 `json:"ram"`
}

type MetricsConsumer struct {
	dlq kafka1.MessageProducer
}

func NewMetricsConsumer(dlq kafka1.MessageProducer) *MetricsConsumer {
	return &MetricsConsumer{dlq: dlq}
}

func (m *MetricsConsumer) Handle(ctx context.Context, msg kafka.Message) error {
	var event MetricsEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return m.dlq.Publish(ctx, "metrics-dlq", msg.Key, msg.Value)
	}

	log.Println("METRICS:", event.CPU, event.RAM)
	return nil
}
