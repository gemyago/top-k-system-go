package aggregation

import (
	"github.com/gemyago/top-k-system-go/internal/di"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewCommands,
		di.ProvideValue(CountersFactory(CountersFactoryFunc(NewCounters))),
		NewCheckPointer,

		// package private deps
		newItemEventsAggregatorModel,
		newItemEventsAggregator,
		newCheckPointerModel,
	)
}
