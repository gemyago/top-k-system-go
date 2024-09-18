package services

import (
	"github.com/segmentio/kafka-go"
)

type ItemEventsKafkaTopicWriter struct {
	*kafka.Writer
}

func NewItemEventsKafkaTopicWriter() *ItemEventsKafkaTopicWriter {
	return &ItemEventsKafkaTopicWriter{
		Writer: &kafka.Writer{
			Topic:                  "item-events",
			AllowAutoTopicCreation: true,                         // TODO: for local mode only
			Addr:                   kafka.TCP("localhost:29092"), // TODO: Configurable
		},
	}
}
