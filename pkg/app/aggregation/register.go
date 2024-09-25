package aggregation

import (
	"github.com/gemyago/top-k-system-go/pkg/di"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewItemEventsAggregator,
		NewItemEventsAggregatorModel,
	)
}
