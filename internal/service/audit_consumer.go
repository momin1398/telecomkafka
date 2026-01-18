package service

import (
	"context"
	"encoding/json"
	"log"

	kafka1 "github.com/momin1398/telecomkafka/internal/kafka"

	"github.com/segmentio/kafka-go"
)

type AuditEvent struct {
	User   string `json:"user"`
	Action string `json:"action"`
}

type AuditConsumer struct {
	dlq kafka1.MessageProducer
}

func NewAuditConsumer(dlq kafka1.MessageProducer) *AuditConsumer {
	return &AuditConsumer{dlq: dlq}
}

func (a *AuditConsumer) Handle(ctx context.Context, msg kafka.Message) error {
	var event AuditEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return a.dlq.Publish(ctx, "audit-dlq", msg.Key, msg.Value)
	}

	log.Println("AUDIT:", event.User, event.Action)
	return nil
}
