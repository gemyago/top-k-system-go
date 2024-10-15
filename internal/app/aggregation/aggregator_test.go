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
				AggregatorModel: newMockItemEventsAggregatorModel(t),
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
			aggregator := newItemEventsAggregator(deps.deps)

			mockModel, _ := deps.deps.AggregatorModel.(*mockItemEventsAggregatorModel)

			offsetBase := rand.Int63n(1000)
			wantItems := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}

			fetchResultChan := make(chan fetchMessageResult)
			mockModel.EXPECT().fetchMessages(ctx, int64(0)).Return(fetchResultChan)

			cnt := newCounters()
			state := aggregationState{
				counters: cnt,
			}

			exit := make(chan error)
			go func() {
				exit <- aggregator.beginAggregating(ctx, state, beginAggregatingOpts{})
			}()
			for i, v := range wantItems {
				mockModel.EXPECT().aggregateItemEvent(int64(i)+offsetBase, &v)
				fetchResultChan <- fetchMessageResult{offset: int64(i) + offsetBase, event: &v}
			}

			cancel()
			gotErr := <-exit
			require.NoError(t, gotErr)
		})
		t.Run("should stop and flush at given offset", func(t *testing.T) {
			deps := newMockDeps(t)
			deps.deps.Verbose = true
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			aggregator := newItemEventsAggregator(deps.deps)

			mockModel, _ := deps.deps.AggregatorModel.(*mockItemEventsAggregatorModel)

			offsetBase := rand.Int63n(1000)
			wantItems := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}

			cnt := newCounters()
			state := aggregationState{
				counters: cnt,
			}

			fetchResultChan := make(chan fetchMessageResult)
			mockModel.EXPECT().fetchMessages(ctx, offsetBase).Return(fetchResultChan)
			mockModel.EXPECT().flushMessages(ctx, cnt)

			exit := make(chan error)
			go func() {
				exit <- aggregator.beginAggregating(ctx, state, beginAggregatingOpts{
					sinceOffset: offsetBase,
					tillOffset:  offsetBase + int64(len(wantItems)-1),
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
			aggregator := newItemEventsAggregator(deps.deps)

			mockModel, _ := deps.deps.AggregatorModel.(*mockItemEventsAggregatorModel)

			fetchResultChan := make(chan fetchMessageResult)
			mockModel.EXPECT().fetchMessages(ctx, int64(0)).Return(fetchResultChan)
			cnt := newCounters()
			state := aggregationState{
				counters: cnt,
			}

			exit := make(chan error)
			go func() {
				exit <- aggregator.beginAggregating(ctx, state, beginAggregatingOpts{})
			}()
			fetchResultChan <- fetchMessageResult{err: errors.New(faker.Word())}

			cancel()
			gotErr := <-exit
			require.NoError(t, gotErr)
		})
		t.Run("should exit when context cancelled", func(t *testing.T) {
			deps := newMockDeps(t)
			ctx, cancel := context.WithCancel(context.Background())
			aggregator := newItemEventsAggregator(deps.deps)

			fetchResultChan := make(chan fetchMessageResult)

			mockModel, _ := deps.deps.AggregatorModel.(*mockItemEventsAggregatorModel)
			mockModel.EXPECT().fetchMessages(ctx, int64(0)).Return(fetchResultChan)
			cnt := newCounters()
			state := aggregationState{
				counters: cnt,
			}

			exit := make(chan error)
			go func() {
				exit <- aggregator.beginAggregating(ctx, state, beginAggregatingOpts{})
			}()
			cancel()
			gotErr := <-exit
			assert.NoError(t, gotErr)
		})
		t.Run("should flush messages on timer", func(t *testing.T) {
			deps := newMockDeps(t)
			ctx, cancel := context.WithCancel(context.Background())
			aggregator := newItemEventsAggregator(deps.deps)

			mockModel, _ := deps.deps.AggregatorModel.(*mockItemEventsAggregatorModel)
			cnt := newCounters()
			state := aggregationState{
				counters: cnt,
			}

			fetchResultChan := make(chan fetchMessageResult)
			mockModel.EXPECT().fetchMessages(ctx, int64(0)).Return(fetchResultChan)
			mockModel.EXPECT().flushMessages(ctx, cnt)

			exit := make(chan error)
			go func() {
				exit <- aggregator.beginAggregating(ctx, state, beginAggregatingOpts{})
			}()
			deps.flushTickerChan <- time.Now()

			cancel()
			gotErr := <-exit
			require.NoError(t, gotErr)
		})
	})
}
