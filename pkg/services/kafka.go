package services

import (
	"context"
	"log/slog"

	"github.com/gemyago/top-k-system-go/pkg/di"
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

	RootLogger *slog.Logger

	KafkaTopic                  string `name:"config.kafka.itemEventsTopic"`
	KafkaAddress                string `name:"config.kafka.address"`
	KafkaAllowAutoTopicCreation bool   `name:"config.kafka.allowAutoTopicCreation"`
}

type ItemEventsKafkaWriterOut struct {
	dig.Out

	Writer          ItemEventsKafkaWriter
	ShutdownHandler di.ProcessShutdownHandler `group:"shutdown-handlers"`
}

func NewItemEventsKafkaWriter(deps ItemEventsKafkaWriterDeps) ItemEventsKafkaWriterOut {
	writer := &kafka.Writer{
		Topic:                  deps.KafkaTopic,
		AllowAutoTopicCreation: deps.KafkaAllowAutoTopicCreation,
		Addr:                   kafka.TCP(deps.KafkaAddress),

		// TODO: This may need some thinking
		Async: true,
	}

	return ItemEventsKafkaWriterOut{
		Writer:          writer,
		ShutdownHandler: di.MakeProcessShutdownHandlerNoContext("Item Events Topic Writer", writer.Close),
	}
}
