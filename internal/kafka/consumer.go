package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, msg kafka.Message) error

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
	}
}

func (c *Consumer) Start(ctx context.Context, handler MessageHandler) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		if err := handler(ctx, msg); err != nil {
			log.Println("processing failed:", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("commit failed:", err)
		}
	}
}
