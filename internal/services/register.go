package services

import (
	"context"
	"time"

	"github.com/gemyago/top-k-system-go/internal/di"
	"github.com/gemyago/top-k-system-go/internal/services/blobstorage"
	"github.com/segmentio/kafka-go"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewTimeProvider,
		NewUUIDGenerator,
		NewItemEventsKafkaReader,
		NewItemEventsKafkaWriter,
		NewShutdownHooks,
		di.ProvideValue(time.NewTicker),
		blobstorage.NewLocalStorage,

		// package private deps
		di.ProvideValue[kafkaLeaderDialer](
			func(
				ctx context.Context, network, addr, topic string, partition int,
			) (kafkaConn, error) { // coverage-ignore // very challenging to test this
				return kafka.DialLeader(ctx, network, addr, topic, partition)
			}),
	)
}
