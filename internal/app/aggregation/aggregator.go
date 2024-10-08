package aggregation

import (
	"context"
	"log/slog"
	"time"

	"github.com/gemyago/top-k-system-go/internal/diag"
	"go.uber.org/dig"
)

type BeginAggregatingOpts struct {
	// TillOffset indicates the offset to aggregate until
	TillOffset int64
}

type ItemEventsAggregator interface {
	BeginAggregating(context context.Context, counters Counters, opts BeginAggregatingOpts) error
}

type ItemEventsAggregatorDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	FlushInterval time.Duration `name:"config.aggregator.flushInterval"`
	Verbose       bool          `name:"config.aggregator.verbose"`

	// app layer
	AggregatorModel ItemEventsAggregatorModel

	// service layer
	TickerFactory func(d time.Duration) *time.Ticker
}

type itemEventsAggregator struct {
	logger *slog.Logger
	ItemEventsAggregatorDeps
}

func (a *itemEventsAggregator) BeginAggregating(
	ctx context.Context,
	counters Counters,
	opts BeginAggregatingOpts,
) error {
	// TODO: Set the offset to start fetching from
	// and keep fetching until the offset provided
	messagesChan := a.AggregatorModel.fetchMessages(ctx)
	flushTimer := a.ItemEventsAggregatorDeps.TickerFactory(a.FlushInterval)
	for {
		select {
		case <-flushTimer.C:
			a.AggregatorModel.flushMessages(ctx, counters)
		case res := <-messagesChan:
			// TODO: Potentially Better error handling here
			if res.err != nil {
				a.logger.ErrorContext(ctx, "failed to fetch message", diag.ErrAttr(res.err))
			} else {
				a.AggregatorModel.aggregateItemEvent(res.offset, res.event)
				if a.Verbose {
					a.logger.DebugContext(ctx, "Item event aggregated",
						slog.String("itemID", res.event.ItemID),
						slog.Int64("offset", res.offset),
					)
				}
				if opts.TillOffset > 0 && res.offset >= opts.TillOffset {
					a.logger.InfoContext(ctx, "Target offset reached. Flushing and stopping aggregation.",
						slog.Int64("offset", res.offset),
						slog.Int64("tillOffset", opts.TillOffset),
					)
					a.AggregatorModel.flushMessages(ctx, counters)
					return nil
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func NewItemEventsAggregator(deps ItemEventsAggregatorDeps) ItemEventsAggregator {
	return &itemEventsAggregator{
		logger:                   deps.RootLogger.WithGroup("item-events-aggregator"),
		ItemEventsAggregatorDeps: deps,
	}
}