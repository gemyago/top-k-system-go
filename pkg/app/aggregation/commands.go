package aggregation

import "context"

type Commands interface {
	// StartAggregator will restore last state and start aggregating
	// events
	StartAggregator(ctx context.Context) error

	// CreateCheckPoint will restore last state, aggregate new events
	// and create a new checkpoint
	CreateCheckPoint(ctx context.Context) error
}
