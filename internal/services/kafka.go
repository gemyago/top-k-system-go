package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/dig"
)

type ItemEventsKafkaWriter struct {
	*kafka.Writer
}

// We may want to remove this once below PR is merged and new version is released:
// https://github.com/segmentio/kafka-go/pull/1341
func (w ItemEventsKafkaWriter) Close() error {
	err := w.Writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close kafka writer: %w", err)
	}
	errorsCount := w.Writer.Stats().Errors
	if errorsCount > 0 {
		return fmt.Errorf("failed to close writer gracefully, %d errors occurred", errorsCount)
	}
	return err
}

type ItemEventsKafkaWriterDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	KafkaTopic                  string        `name:"config.kafka.itemEventsTopic"`
	KafkaAddress                string        `name:"config.kafka.address"`
	KafkaAllowAutoTopicCreation bool          `name:"config.kafka.allowAutoTopicCreation"`
	KafkaWriteTimeout           time.Duration `name:"config.kafka.writeTimeout"`
	KafkaMaxWriteAttempts       int           `name:"config.kafka.maxWriteAttempts"`

	// services
	*ShutdownHooks
}

func NewItemEventsKafkaWriter(deps ItemEventsKafkaWriterDeps) ItemEventsKafkaWriter {
	logger := deps.RootLogger.WithGroup("kafka-writer")
	writer := &kafka.Writer{
		Topic:                  deps.KafkaTopic,
		AllowAutoTopicCreation: deps.KafkaAllowAutoTopicCreation,
		Addr:                   kafka.TCP(deps.KafkaAddress),
		ErrorLogger: kafka.LoggerFunc(func(s string, i ...interface{}) { // coverage-ignore
			// no context here
			logger.ErrorContext(context.Background(), fmt.Sprintf(s, i...))
		}),

		MaxAttempts:  deps.KafkaMaxWriteAttempts,
		WriteTimeout: deps.KafkaWriteTimeout,

		// TODO: This may need some thinking
		Async: true,
	}

	deps.ShutdownHooks.RegisterNoCtx("item-events-writer", writer.Close)

	return ItemEventsKafkaWriter{Writer: writer}
}

type kafkaConn interface {
	ReadLastOffset() (int64, error)
	Close() error
}

type kafkaLeaderDialer func(
	ctx context.Context, network, addr, topic string, partition int,
) (kafkaConn, error)

type ItemEventsKafkaReaderDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	KafkaTopic    string        `name:"config.kafka.itemEventsTopic"`
	KafkaAddress  string        `name:"config.kafka.address"`
	ReaderMaxWait time.Duration `name:"config.kafka.readerMaxWait"`

	// services
	*ShutdownHooks

	// package internal
	KafkaLeaderDialer kafkaLeaderDialer
}

type ItemEventsKafkaReader struct {
	deps ItemEventsKafkaReaderDeps
	*kafka.Reader
}

// ReadLastOffset reads the last offset from the kafka topic. This is going to be an offset
// for the next message produced.
func (r *ItemEventsKafkaReader) ReadLastOffset(ctx context.Context) (int64, error) {
	// TODO: Make partition configurable
	conn, err := r.deps.KafkaLeaderDialer(ctx, "tcp", r.deps.KafkaAddress, r.deps.KafkaTopic, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to dial kafka to read current offset: %w", err)
	}
	defer conn.Close()

	offset, err := conn.ReadLastOffset()
	if err != nil {
		return 0, fmt.Errorf("failed to read last offset: %w", err)
	}
	return offset, nil
}

func NewItemEventsKafkaReader(deps ItemEventsKafkaReaderDeps) *ItemEventsKafkaReader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{deps.KafkaAddress},
		Topic:   deps.KafkaTopic,
		MaxWait: deps.ReaderMaxWait,
		// TODO: Make partition configurable
	})

	deps.ShutdownHooks.RegisterNoCtx("item-events-reader", reader.Close)

	return &ItemEventsKafkaReader{deps: deps, Reader: reader}
}
