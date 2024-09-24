package services

import (
	"context"
	"log/slog"
	"time"

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

type KafkaReader interface {
	Close() error
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	FetchMessage(ctx context.Context) (kafka.Message, error)
	Offset() int64
	SetOffset(offset int64) error
}

type ItemEventsKafkaReader KafkaReader

type ItemEventsKafkaReaderOut struct {
	dig.Out

	Reader          ItemEventsKafkaReader
	ShutdownHandler di.ProcessShutdownHandler `group:"shutdown-handlers"`
}

type ItemEventsKafkaReaderDeps struct {
	dig.In

	RootLogger *slog.Logger

	KafkaTopic    string        `name:"config.kafka.itemEventsTopic"`
	KafkaAddress  string        `name:"config.kafka.address"`
	ReaderMaxWait time.Duration `name:"config.kafka.readerMaxWait"`
}

func NewItemEventsKafkaReader(deps ItemEventsKafkaReaderDeps) ItemEventsKafkaReaderOut {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{deps.KafkaAddress},
		Topic:   deps.KafkaTopic,
		MaxWait: deps.ReaderMaxWait,
	})

	return ItemEventsKafkaReaderOut{
		Reader:          reader,
		ShutdownHandler: di.MakeProcessShutdownHandlerNoContext("Item Events Reader", reader.Close),
	}
}
