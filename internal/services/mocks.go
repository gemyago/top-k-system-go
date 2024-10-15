//go:build !release

package services

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Mock interfaces are used to generate mock implementations of all of the components
// that will be reused elsewhere in a system. This helps to minimize the amount of
// duplicate mock implementations that need to be written.

type mockKafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

var _ mockKafkaWriter = (*kafka.Writer)(nil)

type mockKafkaReader interface {
	Close() error
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	FetchMessage(ctx context.Context) (kafka.Message, error)
	SetOffset(offset int64) error
	ReadLastOffset(ctx context.Context) (int64, error)
}

var _ mockKafkaReader = (*ItemEventsKafkaReader)(nil)
