package ingestion

import (
	"github.com/gemyago/top-k-system-go/internal/di"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewCommands,
	)
}
