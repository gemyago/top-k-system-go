//go:build !release

package aggregation

import "context"

// Mock interfaces are used to generate mock implementations of all of the components
// that will be reused elsewhere in a system. This helps to minimize the amount of
// duplicate mock implementations that need to be written.

type mockCommands interface {
	// StartAggregator will restore last state and start aggregating
	// events
	StartAggregator(ctx context.Context) error

	// CreateCheckPoint will restore last state, aggregate new events
	// and create a new checkpoint
	CreateCheckPoint(ctx context.Context) error
}

var _ mockCommands = (*Commands)(nil)
