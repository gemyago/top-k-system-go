package aggregation

import (
	"context"
)

type ItemEventsAggregatorState struct {
	LastOffset int
}

type ItemEventsAggregator interface {
	RestoreState(context context.Context, state ItemEventsAggregatorState) error
	BeginAggregating(context context.Context) error
}

type itemEventsAggregator struct {
}

func (a *itemEventsAggregator) RestoreState(_ context.Context, _ ItemEventsAggregatorState) error {
	panic("not implemented")
}

func (a *itemEventsAggregator) BeginAggregating(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func NewItemEventsAggregator() ItemEventsAggregator {
	return &itemEventsAggregator{}
}
