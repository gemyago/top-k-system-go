package ingestion

import (
	"github.com/gemyago/top-k-system-go/internal/di"
	"github.com/gemyago/top-k-system-go/internal/services"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		di.ProvideAs[services.ItemEventsKafkaWriter, itemEventsWriter],

		NewCommands,
	)
}
