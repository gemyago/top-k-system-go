package services

import (
	"context"
	"log/slog"
	"time"

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

	// config
	KafkaTopic                  string `name:"config.kafka.itemEventsTopic"`
	KafkaAddress                string `name:"config.kafka.address"`
	KafkaAllowAutoTopicCreation bool   `name:"config.kafka.allowAutoTopicCreation"`

	// services
	ShutdownHooks
}

func NewItemEventsKafkaWriter(deps ItemEventsKafkaWriterDeps) ItemEventsKafkaWriter {
	writer := &kafka.Writer{
		Topic:                  deps.KafkaTopic,
		AllowAutoTopicCreation: deps.KafkaAllowAutoTopicCreation,
		Addr:                   kafka.TCP(deps.KafkaAddress),

		// TODO: This may need some thinking
		Async: true,
	}

	deps.ShutdownHooks.Register(NewShutdownHookNoCtx("Item Events Topic Writer", writer.Close))

	return writer
}

type KafkaReader interface {
	Close() error
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	FetchMessage(ctx context.Context) (kafka.Message, error)
	Offset() int64
	SetOffset(offset int64) error
	Stats() kafka.ReaderStats
	ReadLag(ctx context.Context) (lag int64, err error)
}

type ItemEventsKafkaReader KafkaReader

type ItemEventsKafkaReaderDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	KafkaTopic    string        `name:"config.kafka.itemEventsTopic"`
	KafkaAddress  string        `name:"config.kafka.address"`
	ReaderMaxWait time.Duration `name:"config.kafka.readerMaxWait"`

	// services
	ShutdownHooks
}

func NewItemEventsKafkaReader(deps ItemEventsKafkaReaderDeps) ItemEventsKafkaReader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{deps.KafkaAddress},
		Topic:   deps.KafkaTopic,
		MaxWait: deps.ReaderMaxWait,
	})

	deps.ShutdownHooks.Register(NewShutdownHookNoCtx("Item Events Reader", reader.Close))

	return reader
}
