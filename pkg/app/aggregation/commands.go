package aggregation

import (
	"context"
	"fmt"
	"log/slog"

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

	RootLogger *slog.Logger

	// app layer
	CheckPointer
	ItemEventsAggregator
	CountersFactory
}

type commands struct {
	logger *slog.Logger
	CommandsDeps
}

func (c *commands) StartAggregator(ctx context.Context) error {
	c.logger.InfoContext(ctx, "Restoring counters state")
	counters := c.CountersFactory.NewCounters()
	if err := c.CheckPointer.restoreState(ctx, counters); err != nil {
		return err
	}

	// TODO: Here we need some way to activate counters
	// so then API layer could query them

	c.logger.InfoContext(ctx, "Starting aggregation")
	return c.ItemEventsAggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{})
}

func (c *commands) CreateCheckPoint(ctx context.Context) error {
	counters := c.CountersFactory.NewCounters()
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
		logger:       deps.RootLogger.WithGroup("aggregator.commands"),
		CommandsDeps: deps,
	}
}
