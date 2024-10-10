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
	ReadLastOffset(ctx context.Context) (int64, error)
}

type CommandsDeps struct {
	// all injectable fields must be exported
	// to let dig inject them

	dig.In

	RootLogger *slog.Logger

	// service layer
	ItemEventsReader itemEventsKafkaReader

	// package private components
	ItemEventsAggregator itemEventsAggregator
	CheckPointer         checkPointer
	CountersFactory      countersFactory
}

type Commands struct {
	logger *slog.Logger
	deps   CommandsDeps
}

func (c *Commands) StartAggregator(ctx context.Context) error {
	c.logger.InfoContext(ctx, "Restoring counters state")
	cnt := c.deps.CountersFactory.newCounters()
	if err := c.deps.CheckPointer.restoreState(ctx, cnt); err != nil {
		return fmt.Errorf("failed to restore state while starting aggregator: %w", err)
	}

	// TODO: Here we need some way to activate counters
	// so then API layer could query them

	c.logger.InfoContext(ctx, "Starting aggregation")
	return c.deps.ItemEventsAggregator.beginAggregating(ctx, cnt, beginAggregatingOpts{})
}

func (c *Commands) CreateCheckPoint(ctx context.Context) error {
	ctn := c.deps.CountersFactory.newCounters()

	c.logger.InfoContext(ctx, "Starting creating check point. Restoring last state.")
	if err := c.deps.CheckPointer.restoreState(ctx, ctn); err != nil {
		return fmt.Errorf("failed to restore state while creating check point: %w", err)
	}

	streamTail, err := c.deps.ItemEventsReader.ReadLastOffset(ctx)
	if err != nil {
		return fmt.Errorf("failed to read the lag: %w", err)
	}

	lastOffset := ctn.getLastOffset()
	if lastOffset > 0 {
		// We want to consume starting form the next offset, so doing +1
		if err = c.deps.ItemEventsReader.SetOffset(lastOffset + 1); err != nil {
			return fmt.Errorf("failed to set next offset: %w", err)
		}
	}

	// lag count starts from zero
	if streamTail-lastOffset-1 <= 0 {
		c.logger.InfoContext(ctx,
			"No new messages produced. Checkpoint skipped.",
			slog.Int64("lastOffset", lastOffset),
			slog.Int64("streamTail", streamTail),
		)
		return nil
	}

	c.logger.InfoContext(ctx,
		"Aggregating remaining messages",
		slog.Int64("sinceOffset", lastOffset),
		slog.Int64("streamTail", streamTail),
	)
	if err = c.deps.ItemEventsAggregator.beginAggregating(ctx, ctn, beginAggregatingOpts{
		TillOffset: streamTail,
	}); err != nil {
		return fmt.Errorf("failed to aggregate till offset: %w", err)
	}

	c.logger.InfoContext(ctx, "Producing new state")
	if err = c.deps.CheckPointer.dumpState(ctx, ctn); err != nil {
		return fmt.Errorf("failed to dump state: %w", err)
	}

	c.logger.InfoContext(ctx, "Checkpoint created", slog.Int64("lastOffset", ctn.getLastOffset()))

	return nil
}

func NewCommands(deps CommandsDeps) *Commands {
	return &Commands{
		logger: deps.RootLogger.WithGroup("aggregator.commands"),
		deps:   deps,
	}
}
