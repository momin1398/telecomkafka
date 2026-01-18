// internal/service/telemetry.go
package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/momin1398/telecomkafka/internal/kafka"

	"github.com/momin1398/telecomkafka/internal/model"
)

type TelemetryService struct {
	producer kafka.MessageProducer
	topic    string
}

func NewTelemetryService(p kafka.MessageProducer, topic string) *TelemetryService {
	return &TelemetryService{
		producer: p,
		topic:    topic,
	}
}

func (s *TelemetryService) Process(ctx context.Context, event model.TelemetryEvent) error {
	fmt.Printf("Processing event: %+v\n", event)

	data, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshalling event:", err)
		return err
	}

	err = s.producer.Publish(ctx, s.topic, []byte(event.DeviceID), data)
	if err != nil {
		fmt.Println("Error publishing to Kafka:", err)
		return err
	}

	fmt.Println("Event successfully sent to Kafka")
	return nil
}
