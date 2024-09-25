package aggregation

import (
	"context"
	"fmt"

	"go.uber.org/dig"
)

type Commands interface {
	// StartAggregator will restore last state and start aggregating
	// events
	StartAggregator(ctx context.Context) error

	// CreateCheckPoint will restore last state, aggregate new events
	// and create a new checkpoint
	CreateCheckPoint(ctx context.Context) error
}

type CommandsDeps struct {
	dig.In

	// app layer
	CheckPointer
	ItemEventsAggregator
	CountersFactory
}

type commands struct {
	CommandsDeps
}

func (c *commands) StartAggregator(ctx context.Context) error {
	counters := c.CountersFactory()
	if err := c.CheckPointer.restoreState(ctx, counters); err != nil {
		return err
	}

	// TODO: Here we need some way to activate counters
	// so then API layer could query them

	return c.ItemEventsAggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{})
}

func (c *commands) CreateCheckPoint(ctx context.Context) error {
	counters := c.CountersFactory()
	if err := c.CheckPointer.restoreState(ctx, counters); err != nil {
		return err
	}

	// Get current state
	var currentOffset int64 = 0

	if err := c.ItemEventsAggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{
		TillOffset: currentOffset,
	}); err != nil {
		return fmt.Errorf("failed to aggregate till offset: %w", err)
	}

	if err := c.CheckPointer.dumpState(ctx, counters); err != nil {
		return fmt.Errorf("failed to dump state: %w", err)
	}

	return nil
}

func NewCommands(deps CommandsDeps) Commands {
	return &commands{
		CommandsDeps: deps,
	}
}
