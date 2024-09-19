package services

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type ItemEventsKafkaWriter KafkaWriter

func NewItemEventsKafkaWriter() ItemEventsKafkaWriter {
	// TODO: Need to close on shutdown to make sure pending events got flushed
	return &kafka.Writer{
		Topic:                  "item-events",
		AllowAutoTopicCreation: true,                         // TODO: for local mode only
		Addr:                   kafka.TCP("localhost:29092"), // TODO: Configurable

		// TODO: This may need some thinking
		Async: true,
	}
}
