package aggregation

import (
	"context"
	"log/slog"
	"time"

	"github.com/gemyago/top-k-system-go/internal/diag"
	"go.uber.org/dig"
)

type aggregationState struct {
	counters     counters
	allTimeItems topKItems
}

type beginAggregatingOpts struct {
	sinceOffset int64

	// tillOffset indicates the offset to aggregate until
	tillOffset int64
}

type itemEventsAggregator interface {
	beginAggregating(context context.Context, state aggregationState, opts beginAggregatingOpts) error
}

type ItemEventsAggregatorDeps struct {
	// all injectable fields must be exported
	// to let dig inject them

	dig.In

	RootLogger *slog.Logger

	// config
	FlushInterval    time.Duration `name:"config.aggregator.flushInterval"`
	Verbose          bool          `name:"config.aggregator.verbose"`
	ItemEventLogRate int64         `name:"config.aggregator.itemEventLogRate"`

	// service layer
	TickerFactory func(d time.Duration) *time.Ticker

	// package private components
	AggregatorModel itemEventsAggregatorModel
}

type itemEventsAggregatorImpl struct {
	logger *slog.Logger
	ItemEventsAggregatorDeps
}

func (a *itemEventsAggregatorImpl) beginAggregating(
	ctx context.Context,
	state aggregationState,
	opts beginAggregatingOpts,
) error {
	messagesChan := a.AggregatorModel.fetchMessages(ctx, opts.sinceOffset)
	flushTimer := a.ItemEventsAggregatorDeps.TickerFactory(a.FlushInterval)
	for {
		select {
		case <-flushTimer.C:
			a.AggregatorModel.flushMessages(ctx, state)
		case res := <-messagesChan:
			// TODO: Potentially Better error handling here
			if res.err != nil {
				a.logger.ErrorContext(ctx, "failed to fetch message", diag.ErrAttr(res.err))
			} else {
				a.AggregatorModel.aggregateItemEvent(res.offset, res.event)
				shouldLog := a.Verbose || (a.ItemEventLogRate > 0 && res.offset%a.ItemEventLogRate == 0)
				if shouldLog {
					a.logger.DebugContext(ctx, "Item event aggregated",
						slog.String("itemID", res.event.ItemID),
						slog.Int64("offset", res.offset),
					)
				}
				if opts.tillOffset > 0 && res.offset >= opts.tillOffset {
					a.logger.InfoContext(ctx, "Target offset reached. Flushing and stopping aggregation.",
						slog.Int64("offset", res.offset),
						slog.Int64("tillOffset", opts.tillOffset),
					)
					a.AggregatorModel.flushMessages(ctx, state)
					return nil
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func newItemEventsAggregator(deps ItemEventsAggregatorDeps) itemEventsAggregator {
	return &itemEventsAggregatorImpl{
		logger:                   deps.RootLogger.WithGroup("item-events-aggregator"),
		ItemEventsAggregatorDeps: deps,
	}
}
