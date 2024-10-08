package aggregation

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/gemyago/top-k-system-go/internal/app/models"
	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemEventsAggregator(t *testing.T) {
	type itemEventsAggregatorMockDeps struct {
		deps            ItemEventsAggregatorDeps
		flushTickerChan chan time.Time
	}

	newMockDeps := func(t *testing.T) itemEventsAggregatorMockDeps {
		flushTickerChan := make(chan time.Time)
		flushTicker := &time.Ticker{C: flushTickerChan}
		flushInterval := time.Duration(rand.Int63n(1000))
		return itemEventsAggregatorMockDeps{
			flushTickerChan: flushTickerChan,
			deps: ItemEventsAggregatorDeps{
				RootLogger:      diag.RootTestLogger(),
				AggregatorModel: NewMockItemEventsAggregatorModel(t),
				FlushInterval:   flushInterval,
				TickerFactory: func(d time.Duration) *time.Ticker {
					assert.Equal(t, flushInterval, d)
					return flushTicker
				},
			},
		}
	}

	t.Run("BeginAggregating", func(t *testing.T) {
		t.Run("should aggregate messages", func(t *testing.T) {
			deps := newMockDeps(t)
			deps.deps.Verbose = true
			ctx, cancel := context.WithCancel(context.Background())
			aggregator := NewItemEventsAggregator(deps.deps)

			mockModel, _ := deps.deps.AggregatorModel.(*MockItemEventsAggregatorModel)

			offsetBase := rand.Int63n(1000)
			wantItems := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}

			fetchResultChan := make(chan fetchMessageResult)
			mockModel.EXPECT().fetchMessages(ctx).Return(fetchResultChan)

			counters := NewCounters()

			exit := make(chan error)
			go func() {
				exit <- aggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{})
			}()
			for i, v := range wantItems {
				mockModel.EXPECT().aggregateItemEvent(int64(i)+offsetBase, &v)
				fetchResultChan <- fetchMessageResult{offset: int64(i) + offsetBase, event: &v}
			}

			cancel()
			gotErr := <-exit
			require.NoError(t, gotErr)
			mockModel.AssertExpectations(t)
		})
		t.Run("should stop and flush at given offset", func(t *testing.T) {
			deps := newMockDeps(t)
			deps.deps.Verbose = true
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			aggregator := NewItemEventsAggregator(deps.deps)

			mockModel, _ := deps.deps.AggregatorModel.(*MockItemEventsAggregatorModel)

			offsetBase := rand.Int63n(1000)
			wantItems := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}

			counters := NewCounters()
			fetchResultChan := make(chan fetchMessageResult)
			mockModel.EXPECT().fetchMessages(ctx).Return(fetchResultChan)
			mockModel.EXPECT().flushMessages(ctx, counters)

			exit := make(chan error)
			go func() {
				exit <- aggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{
					TillOffset: offsetBase + int64(len(wantItems)-1),
				})
			}()
			for i, v := range wantItems {
				mockModel.EXPECT().aggregateItemEvent(int64(i)+offsetBase, &v)
				fetchResultChan <- fetchMessageResult{offset: int64(i) + offsetBase, event: &v}
			}
			gotErr := <-exit
			require.NoError(t, gotErr)
		})
		t.Run("should handle errors when fetch messages", func(t *testing.T) {
			deps := newMockDeps(t)
			ctx, cancel := context.WithCancel(context.Background())
			aggregator := NewItemEventsAggregator(deps.deps)

			mockModel, _ := deps.deps.AggregatorModel.(*MockItemEventsAggregatorModel)

			fetchResultChan := make(chan fetchMessageResult)
			mockModel.EXPECT().fetchMessages(ctx).Return(fetchResultChan)
			counters := NewCounters()

			exit := make(chan error)
			go func() {
				exit <- aggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{})
			}()
			fetchResultChan <- fetchMessageResult{err: errors.New(faker.Word())}

			cancel()
			gotErr := <-exit
			require.NoError(t, gotErr)
			mockModel.AssertExpectations(t)
		})
		t.Run("should exit when context cancelled", func(t *testing.T) {
			deps := newMockDeps(t)
			ctx, cancel := context.WithCancel(context.Background())
			aggregator := NewItemEventsAggregator(deps.deps)

			fetchResultChan := make(chan fetchMessageResult)

			mockModel, _ := deps.deps.AggregatorModel.(*MockItemEventsAggregatorModel)
			mockModel.EXPECT().fetchMessages(ctx).Return(fetchResultChan)
			counters := NewCounters()

			exit := make(chan error)
			go func() {
				exit <- aggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{})
			}()
			cancel()
			gotErr := <-exit
			assert.NoError(t, gotErr)
		})
		t.Run("should flush messages on timer", func(t *testing.T) {
			deps := newMockDeps(t)
			ctx, cancel := context.WithCancel(context.Background())
			aggregator := NewItemEventsAggregator(deps.deps)

			mockModel, _ := deps.deps.AggregatorModel.(*MockItemEventsAggregatorModel)
			counters := NewCounters()

			fetchResultChan := make(chan fetchMessageResult)
			mockModel.EXPECT().fetchMessages(ctx).Return(fetchResultChan)
			mockModel.EXPECT().flushMessages(ctx, counters)

			exit := make(chan error)
			go func() {
				exit <- aggregator.BeginAggregating(ctx, counters, BeginAggregatingOpts{})
			}()
			deps.flushTickerChan <- time.Now()

			cancel()
			gotErr := <-exit
			require.NoError(t, gotErr)
			mockModel.AssertExpectations(t)
		})
	})
}