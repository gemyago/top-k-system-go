package aggregation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"go.uber.org/dig"
)

type itemEventsKafkaReader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	SetOffset(offset int64) error
	ReadLag(ctx context.Context) (lag int64, err error)
}

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
	countersFactory

	// service layer
	ItemEventsReader itemEventsKafkaReader

	// package private components
	itemEventsAggregator
	checkPointer
}

type commands struct {
	logger *slog.Logger
	CommandsDeps
}

func (c *commands) StartAggregator(ctx context.Context) error {
	c.logger.InfoContext(ctx, "Restoring counters state")
	cnt := c.countersFactory.newCounters()
	if err := c.checkPointer.restoreState(ctx, cnt); err != nil {
		return fmt.Errorf("failed to restore state while starting aggregator: %w", err)
	}

	// TODO: Here we need some way to activate counters
	// so then API layer could query them

	c.logger.InfoContext(ctx, "Starting aggregation")
	return c.itemEventsAggregator.beginAggregating(ctx, cnt, beginAggregatingOpts{})
}

func (c *commands) CreateCheckPoint(ctx context.Context) error {
	ctn := c.countersFactory.newCounters()

	c.logger.InfoContext(ctx, "Starting creating check point. Restoring last state.")
	if err := c.checkPointer.restoreState(ctx, ctn); err != nil {
		return fmt.Errorf("failed to restore state while creating check point: %w", err)
	}

	lag, err := c.ItemEventsReader.ReadLag(ctx)
	if err != nil {
		return fmt.Errorf("failed to read the lag: %w", err)
	}

	lastOffset := ctn.getLastOffset()
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
		slog.Int64("sinceOffset", ctn.getLastOffset()),
		slog.Int64("tillOffset", tillOffset),
	)
	if err = c.itemEventsAggregator.beginAggregating(ctx, ctn, beginAggregatingOpts{
		TillOffset: tillOffset,
	}); err != nil {
		return fmt.Errorf("failed to aggregate till offset: %w", err)
	}

	c.logger.InfoContext(ctx, "Producing new state")
	if err = c.checkPointer.dumpState(ctx, ctn); err != nil {
		return fmt.Errorf("failed to dump state: %w", err)
	}

	c.logger.InfoContext(ctx, "Checkpoint created", slog.Int64("lastOffset", ctn.getLastOffset()))

	return nil
}

func NewCommands(deps CommandsDeps) Commands {
	return &commands{
		logger:       deps.RootLogger.WithGroup("aggregator.commands"),
		CommandsDeps: deps,
	}
}
