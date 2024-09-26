package aggregation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gemyago/top-k-system-go/pkg/services"
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

	// service layer
	ItemEventsReader services.ItemEventsKafkaReader
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

	c.logger.InfoContext(ctx, "Starting creating check point. Restoring last state.")
	if err := c.CheckPointer.restoreState(ctx, counters); err != nil {
		return err
	}

	lag, err := c.ItemEventsReader.ReadLag(ctx)
	if err != nil {
		return fmt.Errorf("failed to read the lag: %w", err)
	}

	lastOffset := counters.getLastOffset()
	if lastOffset > 0 {
		// We want to consume starting form the next offset, so doing +1
		if err = c.ItemEventsReader.SetOffset(lastOffset + 1); err != nil {
			return fmt.Errorf("failed to set next offset: %w", err)
		}
	}

	// lag count starts from zero
	if lag-lastOffset-1 <= 0 {
		c.logger.InfoContext(ctx,
			"No new messages produced. Checkpoint skipped.",
			slog.Int64("lastOffset", lastOffset),
			slog.Int64("lag", lag),
		)
		return nil
	}

	// We assume we didn't consume anything yet and the lag is exactly the
	// tail of the stream
	tillOffset := lag - 1

	c.logger.InfoContext(ctx,
		"Aggregating remaining messages",
		slog.Int64("sinceOffset", counters.getLastOffset()),
		slog.Int64("tillOffset", tillOffset),
	)
	if err = c.ItemEventsAggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{
		TillOffset: tillOffset,
	}); err != nil {
		return fmt.Errorf("failed to aggregate till offset: %w", err)
	}

	c.logger.InfoContext(ctx, "Producing new state")
	if err = c.CheckPointer.dumpState(ctx, counters); err != nil {
		return fmt.Errorf("failed to dump state: %w", err)
	}

	c.logger.InfoContext(ctx, "Checkpoint created", slog.Int64("lastOffset", counters.getLastOffset()))

	return nil
}

func NewCommands(deps CommandsDeps) Commands {
	return &commands{
		logger:       deps.RootLogger.WithGroup("aggregator.commands"),
		CommandsDeps: deps,
	}
}
