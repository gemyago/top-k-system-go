package aggregation

import (
	"context"
	"encoding/json"
	"fmt"

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

	// app layer
	Counters

	// service layer
	ItemEventsReader services.ItemEventsKafkaReader
}

type itemEventsAggregatorModel struct {
	lastAggregatedOffset int64
	aggregatedItems      map[string]int64

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

func (m *itemEventsAggregatorModel) flushMessages(_ context.Context) error {
	panic("not implemented")
}

func NewItemEventsAggregatorModel(
	deps ItemEventsAggregatorModelDeps,
) ItemEventsAggregatorModel {
	return &itemEventsAggregatorModel{
		aggregatedItems: make(map[string]int64),
		deps:            deps,
	}
}
