package aggregation

import (
	"context"

	"github.com/gemyago/top-k-system-go/pkg/app/models"
)

type fetchMessageResult struct {
	event  *models.ItemEvent
	offset int
	err    error
}

type ItemEventsAggregatorModel interface {
	aggregateItemEvent(evt *models.ItemEvent)
	fetchMessages(ctx context.Context) <-chan fetchMessageResult
	flushMessages(ctx context.Context) error
}
