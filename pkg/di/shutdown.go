package di

import "context"

// ProcessShutdownHandler is used mostly to register
// shutdown handlers that should perform cleanup tasks on
// process shutdown.
type ProcessShutdownHandler interface {
	Shutdown(ctx context.Context) error
}

type ProcessShutdownHandlerFunc func(ctx context.Context) error

func (h ProcessShutdownHandlerFunc) Shutdown(ctx context.Context) error {
	return h(ctx)
}
