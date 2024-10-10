package aggregation

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gemyago/top-k-system-go/internal/app/models"
	"go.uber.org/dig"
)

type fetchMessageResult struct {
	event  *models.ItemEvent
	offset int64
	err    error
}

type ItemEventsAggregatorModelDeps struct {
	// all injectable fields must be exported
	// to let dig inject them

	dig.In

	RootLogger *slog.Logger

	// config
	Verbose bool `name:"config.aggregator.verbose"`

	// service layer
	ItemEventsReader itemEventsKafkaReader
}

type itemEventsAggregatorModel interface {
	aggregateItemEvent(offset int64, evt *models.ItemEvent)
	flushMessages(ctx context.Context, counters counters)
	fetchMessages(ctx context.Context, fromOffset int64) <-chan fetchMessageResult
}

type itemEventsAggregatorModelImpl struct {
	lastAggregatedOffset int64
	aggregatedItems      map[string]int64
	logger               *slog.Logger

	deps ItemEventsAggregatorModelDeps
}

// aggregateItemEvent method is not thread safe, should be only called from a same
// goroutine as flushMessages.
func (m *itemEventsAggregatorModelImpl) aggregateItemEvent(offset int64, evt *models.ItemEvent) {
	m.lastAggregatedOffset = offset
	curVal := m.aggregatedItems[evt.ItemID]
	m.aggregatedItems[evt.ItemID] = curVal + 1
}

// flushMessages method is not thread safe, should be only called from a same
// goroutine as aggregateItemEvent.
func (m *itemEventsAggregatorModelImpl) flushMessages(ctx context.Context, counters counters) {
	m.logger.DebugContext(ctx, "Flushing aggregated messages")
	counters.updateItemsCount(m.lastAggregatedOffset, m.aggregatedItems)
	clear(m.aggregatedItems)
}

func (m *itemEventsAggregatorModelImpl) fetchMessages(ctx context.Context, fromOffset int64) <-chan fetchMessageResult {
	resultsChan := make(chan fetchMessageResult)
	if err := m.deps.ItemEventsReader.SetOffset(fromOffset); err != nil {
		resultsChan <- fetchMessageResult{err: fmt.Errorf("failed to set offset: %w", err)}
		close(resultsChan)
		return resultsChan
	}

	go func() {
		for {
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

func newItemEventsAggregatorModel(
	deps ItemEventsAggregatorModelDeps,
) itemEventsAggregatorModel {
	return &itemEventsAggregatorModelImpl{
		logger:          deps.RootLogger.WithGroup("item-events-aggregator-model"),
		aggregatedItems: make(map[string]int64),
		deps:            deps,
	}
}
