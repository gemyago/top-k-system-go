package services

import (
	"time"

	"github.com/gemyago/top-k-system-go/pkg/di"
	"github.com/gemyago/top-k-system-go/pkg/services/blobstorage"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewTimeProvider,
		NewItemEventsKafkaReader,
		NewItemEventsKafkaWriter,
		di.ProvideValue(time.NewTicker),
		blobstorage.NewLocalStorage,
	)
}
