package aggregation

import (
	"context"
	"log/slog"
	"time"

	"github.com/gemyago/top-k-system-go/pkg/diag"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"go.uber.org/dig"
)

type ItemEventsAggregatorState struct {
	LastOffset int
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

	// app layer
	AggregatorModel ItemEventsAggregatorModel
	Counters

	// service layer
	ItemEventsReader services.ItemEventsKafkaReader
	TickerFactory    func(d time.Duration) *time.Ticker
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
			// TODO: Potentially Better error handling here
			// Like if 3 events in a row then panic or something
			if err := a.AggregatorModel.flushMessages(ctx); err != nil {
				a.logger.ErrorContext(ctx, "failed to flush aggregated messages", diag.ErrAttr(err))
			}
		case res := <-messagesChan:
			// TODO: Potentially Better error handling here
			if res.err != nil {
				a.logger.ErrorContext(ctx, "failed to fetch message", diag.ErrAttr(res.err))
			} else {
				a.AggregatorModel.aggregateItemEvent(res.event)
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
