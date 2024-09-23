package aggregation

import (
	"context"

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
}

// aggregateItemEvent method is not thread safe, should be only called from a single
// goroutine.
func (m *itemEventsAggregatorModel) aggregateItemEvent(offset int64, evt *models.ItemEvent) {
	m.lastAggregatedOffset = offset
	curVal := m.aggregatedItems[evt.ItemID]
	m.aggregatedItems[evt.ItemID] = curVal + 1
}

func (m *itemEventsAggregatorModel) fetchMessages(ctx context.Context) <-chan fetchMessageResult {
	panic("not implemented")
}

func (m *itemEventsAggregatorModel) flushMessages(ctx context.Context) error {
	panic("not implemented")
}

func NewItemEventsAggregatorModel(
	deps ItemEventsAggregatorDeps,
) ItemEventsAggregatorModel {
	return &itemEventsAggregatorModel{
		aggregatedItems: make(map[string]int64),
	}
}
