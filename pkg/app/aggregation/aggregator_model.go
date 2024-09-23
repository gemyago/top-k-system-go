package aggregation

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"go.uber.org/dig"
)

type fetchMessageResult struct {
	event  *models.ItemEvent
	offset int64
	err    error
}

type ItemEventsAggregatorModel interface {
	aggregateItemEvent(offset int64, evt *models.ItemEvent)
	fetchMessages(ctx context.Context) <-chan fetchMessageResult
	flushMessages(ctx context.Context) error
}

type ItemEventsAggregatorModelDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	Verbose bool `name:"config.aggregator.verbose"`

	// app layer
	Counters

	// service layer
	ItemEventsReader services.ItemEventsKafkaReader
}

type itemEventsAggregatorModel struct {
	lastAggregatedOffset int64
	aggregatedItems      map[string]int64
	logger               *slog.Logger

	deps ItemEventsAggregatorModelDeps
}

// aggregateItemEvent method is not thread safe, should be only called from a single
// goroutine.
func (m *itemEventsAggregatorModel) aggregateItemEvent(offset int64, evt *models.ItemEvent) {
	m.lastAggregatedOffset = offset
	curVal := m.aggregatedItems[evt.ItemID]
	m.aggregatedItems[evt.ItemID] = curVal + 1
}

func (m *itemEventsAggregatorModel) fetchMessages(ctx context.Context) <-chan fetchMessageResult {
	resultsChan := make(chan fetchMessageResult)
	go func() {
		for {
			if m.deps.Verbose {
				m.logger.DebugContext(ctx, "Waiting for next item message")
			}
			msg, err := m.deps.ItemEventsReader.FetchMessage(ctx)
			if err != nil {
				resultsChan <- fetchMessageResult{err: fmt.Errorf("failed to fetch messages: %w", err)}
				// TODO: If EOF just stop the loop
				// review usage to make sure it will not break anything
			} else {
				var itemEvent models.ItemEvent
				if err = json.Unmarshal(msg.Value, &itemEvent); err != nil {
					resultsChan <- fetchMessageResult{err: fmt.Errorf("failed to unmarshal message: %w", err)}
				}
				resultsChan <- fetchMessageResult{
					event:  &itemEvent,
					offset: msg.Offset,
				}
			}
		}
	}()
	return resultsChan
}

func (m *itemEventsAggregatorModel) flushMessages(ctx context.Context) error {
	m.logger.DebugContext(ctx, "Flushing aggregated messages")
	return nil
}

func NewItemEventsAggregatorModel(
	deps ItemEventsAggregatorModelDeps,
) ItemEventsAggregatorModel {
	return &itemEventsAggregatorModel{
		logger:          deps.RootLogger.WithGroup("item-events-aggregator-model"),
		aggregatedItems: make(map[string]int64),
		deps:            deps,
	}
}
