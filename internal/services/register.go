package services

import (
	"time"

	"github.com/gemyago/top-k-system-go/internal/di"
	"github.com/gemyago/top-k-system-go/internal/services/blobstorage"
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
	)
}
