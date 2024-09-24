package aggregation

import (
	"context"
	"log/slog"
	"time"

	"github.com/gemyago/top-k-system-go/pkg/diag"
	"go.uber.org/dig"
)

type ItemEventsAggregatorState struct {
	LastOffset int64
}

type ItemEventsAggregator interface {
	RestoreState(context context.Context, state ItemEventsAggregatorState) error
	BeginAggregating(context context.Context) error
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

func (a *itemEventsAggregator) RestoreState(_ context.Context, _ ItemEventsAggregatorState) error {
	panic("not implemented")
}

func (a *itemEventsAggregator) BeginAggregating(ctx context.Context) error {
	messagesChan := a.AggregatorModel.fetchMessages(ctx)
	flushTimer := a.ItemEventsAggregatorDeps.TickerFactory(a.FlushInterval)
	for {
		select {
		case <-flushTimer.C:
			a.AggregatorModel.flushMessages(ctx)
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
