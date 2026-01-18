package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"

	kafka1 "github.com/momin1398/telecomkafka/internal/kafka"
)

type AlarmEvent struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type AlarmConsumer struct {
	dlq kafka1.MessageProducer
}

func NewAlarmConsumer(dlq kafka1.MessageProducer) *AlarmConsumer {
	return &AlarmConsumer{dlq: dlq}
}

func (a *AlarmConsumer) Handle(ctx context.Context, msg kafka.Message) error {
	var event AlarmEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return a.dlq.Publish(ctx, "alarm-dlq", msg.Key, msg.Value)
	}

	if event.Severity == "" {
		return errors.New("missing severity")
	}

	log.Println("ALARM:", event.Severity, event.Message)
	return nil
}
