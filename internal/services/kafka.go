package services

import (
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/dig"
)

type ItemEventsKafkaWriter struct {
	*kafka.Writer
}

type ItemEventsKafkaWriterDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	KafkaTopic                  string `name:"config.kafka.itemEventsTopic"`
	KafkaAddress                string `name:"config.kafka.address"`
	KafkaAllowAutoTopicCreation bool   `name:"config.kafka.allowAutoTopicCreation"`

	// services
	*ShutdownHooks
}

func NewItemEventsKafkaWriter(deps ItemEventsKafkaWriterDeps) ItemEventsKafkaWriter {
	writer := &kafka.Writer{
		Topic:                  deps.KafkaTopic,
		AllowAutoTopicCreation: deps.KafkaAllowAutoTopicCreation,
		Addr:                   kafka.TCP(deps.KafkaAddress),

		// TODO: This may need some thinking
		Async: true,
	}

	deps.ShutdownHooks.RegisterNoCtx("item-events-writer", writer.Close)

	return ItemEventsKafkaWriter{Writer: writer}
}

type ItemEventsKafkaReader struct {
	*kafka.Reader
}

type ItemEventsKafkaReaderDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	KafkaTopic    string        `name:"config.kafka.itemEventsTopic"`
	KafkaAddress  string        `name:"config.kafka.address"`
	ReaderMaxWait time.Duration `name:"config.kafka.readerMaxWait"`

	// services
	*ShutdownHooks
}

func NewItemEventsKafkaReader(deps ItemEventsKafkaReaderDeps) ItemEventsKafkaReader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{deps.KafkaAddress},
		Topic:   deps.KafkaTopic,
		MaxWait: deps.ReaderMaxWait,
	})

	deps.ShutdownHooks.RegisterNoCtx("item-events-reader", reader.Close)

	return ItemEventsKafkaReader{Reader: reader}
}
