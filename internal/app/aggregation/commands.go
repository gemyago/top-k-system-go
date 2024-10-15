package aggregation

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/samber/lo"
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
	TopKItemsFactory     topKItemsFactory
	AggregationState     aggregationState
}

type Commands struct {
	logger *slog.Logger
	deps   CommandsDeps
}

func (c *Commands) StartAggregator(ctx context.Context) error {
	c.logger.DebugContext(ctx, "Restoring counters state")
	startedAt := time.Now()
	if err := c.deps.CheckPointer.restoreState(ctx, c.deps.AggregationState); err != nil {
		return fmt.Errorf("failed to restore state while starting aggregator: %w", err)
	}

	counters := c.deps.AggregationState.counters
	c.logger.InfoContext(ctx, "Counters state restored",
		slog.Int("totalItemsCount", len(counters.getItemsCounters())),
		slog.Int64("lastOffset", counters.getLastOffset()),
		slog.Duration("restorationDuration", time.Since(startedAt)),
	)
	lastOffset := counters.getLastOffset()
	sinceOffset := lo.If(lastOffset == 0, int64(0)).Else(lastOffset + 1)
	c.logger.InfoContext(ctx,
		"Starting aggregation",
		slog.Int64("sinceOffset", sinceOffset),
	)
	return c.deps.ItemEventsAggregator.beginAggregating(ctx, c.deps.AggregationState, beginAggregatingOpts{
		sinceOffset: sinceOffset,
	})
}

func (c *Commands) CreateCheckPoint(ctx context.Context) error {
	ctn := c.deps.CountersFactory.newCounters()
	allTimesItems := c.deps.TopKItemsFactory.newTopKItems(topKMaxItemsSize)
	state := aggregationState{
		counters:     ctn,
		allTimeItems: allTimesItems,
	}

	c.logger.InfoContext(ctx, "Starting creating check point. Restoring last state.")
	if err := c.deps.CheckPointer.restoreState(ctx, state); err != nil {
		return fmt.Errorf("failed to restore state while creating check point: %w", err)
	}

	streamTail, err := c.deps.ItemEventsReader.ReadLastOffset(ctx)
	if err != nil {
		return fmt.Errorf("failed to read the lag: %w", err)
	}

	lastOffset := ctn.getLastOffset()

	// the streamTail will have a next offset
	if streamTail-lastOffset-1 <= 0 {
		c.logger.InfoContext(ctx,
			"No new messages produced. Checkpoint skipped.",
			slog.Int64("lastOffset", lastOffset),
			slog.Int64("streamTail", streamTail),
		)
		return nil
	}

	sinceOffset := lo.If(lastOffset == 0, int64(0)).Else(lastOffset + 1)
	c.logger.InfoContext(ctx,
		"Aggregating remaining messages",
		slog.Int64("sinceOffset", sinceOffset),
		slog.Int64("streamTail", streamTail),
	)
	if err = c.deps.ItemEventsAggregator.beginAggregating(ctx, state, beginAggregatingOpts{
		sinceOffset: sinceOffset,
		tillOffset:  streamTail - 1,
	}); err != nil {
		return fmt.Errorf("failed to aggregate till offset: %w", err)
	}

	c.logger.InfoContext(ctx, "Producing new state")
	if err = c.deps.CheckPointer.dumpState(ctx, state); err != nil {
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
