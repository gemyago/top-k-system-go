package services

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/dig"
)

type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type ItemEventsKafkaWriter KafkaWriter

type ItemEventsKafkaWriterDeps struct {
	dig.In

	KafkaTopic                  string `name:"config.kafka.itemEventsTopic"`
	KafkaAddress                string `name:"config.kafka.address"`
	KafkaAllowAutoTopicCreation bool   `name:"config.kafka.allowAutoTopicCreation"`
}

func NewItemEventsKafkaWriter(deps ItemEventsKafkaWriterDeps) ItemEventsKafkaWriter {
	// TODO: Need to close on shutdown to make sure pending events got flushed
	return &kafka.Writer{
		Topic:                  deps.KafkaTopic,
		AllowAutoTopicCreation: deps.KafkaAllowAutoTopicCreation,
		Addr:                   kafka.TCP(deps.KafkaAddress),

		// TODO: This may need some thinking
		Async: true,
	}
}
