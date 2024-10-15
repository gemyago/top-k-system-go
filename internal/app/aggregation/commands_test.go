package aggregation

import (
	"context"
	"errors"
	"math/rand/v2"
	"testing"

	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCommands(t *testing.T) {
	newMockDeps := func(t *testing.T) CommandsDeps {
		return CommandsDeps{
			RootLogger:           diag.RootTestLogger(),
			CheckPointer:         newMockCheckPointer(t),
			ItemEventsAggregator: newMockItemEventsAggregator(t),
			ItemEventsReader:     services.NewMockKafkaReader(t),
			CountersFactory:      newMockCountersFactory(t),
			TopKItemsFactory:     newMockTopKItemsFactory(t),
			AggregationState: aggregationState{
				counters:     newMockCounters(t),
				allTimeItems: newMockTopKItems(t),
			},
		}
	}

	t.Run("StartAggregator", func(t *testing.T) {
		t.Run("should restore state and start aggregating", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()
			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, mockDeps.AggregationState).Return(nil)

			lastOffset := rand.Int64N(100)
			mockCounters, _ := mockDeps.AggregationState.counters.(*mockCounters)
			mockCounters.EXPECT().getLastOffset().Return(lastOffset)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*mockItemEventsAggregator)
			aggregator.EXPECT().
				beginAggregating(ctx, mockDeps.AggregationState, beginAggregatingOpts{
					sinceOffset: lastOffset + 1,
				}).
				Return(nil)

			require.NoError(t, commands.StartAggregator(ctx))
		})

		t.Run("should start aggregating from the beginning if no state", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, mock.Anything).Return(nil)

			mockCounters, _ := mockDeps.AggregationState.counters.(*mockCounters)
			mockCounters.EXPECT().getLastOffset().Return(0)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*mockItemEventsAggregator)
			aggregator.EXPECT().
				beginAggregating(ctx, mockDeps.AggregationState, beginAggregatingOpts{
					sinceOffset: 0,
				}).
				Return(nil)

			require.NoError(t, commands.StartAggregator(ctx))
		})

		t.Run("should return error if restore state failed", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			wantErr := errors.New(faker.Sentence())
			checkPointer.EXPECT().restoreState(ctx, mock.Anything).Return(wantErr)

			require.ErrorIs(t, commands.StartAggregator(ctx), wantErr)
		})
	})

	t.Run("CreateCheckPoint", func(t *testing.T) {
		t.Run("should restore state and aggregate till tail of the queue", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)
			wantCounters.EXPECT().getLastOffset().Return(0)

			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			mockAllTimesItems := newMockTopKItems(t)
			topKItemsFactory, _ := mockDeps.TopKItemsFactory.(*mockTopKItemsFactory)
			topKItemsFactory.EXPECT().newTopKItems(topKMaxItemsSize).Return(mockAllTimesItems)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, aggregationState{
				counters:     wantCounters,
				allTimeItems: mockAllTimesItems,
			}).Return(nil)

			state := aggregationState{
				counters:     wantCounters,
				allTimeItems: mockAllTimesItems,
			}

			wantTail := rand.Int64()
			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLastOffset(ctx).Return(wantTail, nil)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*mockItemEventsAggregator)
			aggregator.EXPECT().
				beginAggregating(ctx, state, beginAggregatingOpts{
					sinceOffset: 0,
					tillOffset:  wantTail - 1,
				}).
				Return(nil)

			checkPointer.EXPECT().dumpState(ctx, state).Return(nil)

			require.NoError(t, commands.CreateCheckPoint(ctx))
		})
		t.Run("should set the restored offset of the reader", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)

			lastOffset := rand.Int64()
			wantCounters.EXPECT().getLastOffset().Return(lastOffset)

			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			mockAllTimesItems := newMockTopKItems(t)
			topKItemsFactory, _ := mockDeps.TopKItemsFactory.(*mockTopKItemsFactory)
			topKItemsFactory.EXPECT().newTopKItems(topKMaxItemsSize).Return(mockAllTimesItems)

			state := aggregationState{
				counters:     wantCounters,
				allTimeItems: mockAllTimesItems,
			}

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, mock.Anything).Return(nil)

			wantTail := lastOffset + 100
			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLastOffset(ctx).Return(wantTail, nil)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*mockItemEventsAggregator)
			aggregator.EXPECT().
				beginAggregating(ctx, state, beginAggregatingOpts{
					sinceOffset: lastOffset + 1,
					tillOffset:  wantTail - 1,
				}).
				Return(nil)

			checkPointer.EXPECT().dumpState(ctx, mock.Anything).Return(nil)

			require.NoError(t, commands.CreateCheckPoint(ctx))
		})
		t.Run("should not aggregate if no new messages", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)

			wantTail := rand.Int64()
			wantCounters.EXPECT().getLastOffset().Return(wantTail)

			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			mockAllTimesItems := newMockTopKItems(t)
			topKItemsFactory, _ := mockDeps.TopKItemsFactory.(*mockTopKItemsFactory)
			topKItemsFactory.EXPECT().newTopKItems(topKMaxItemsSize).Return(mockAllTimesItems)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, mock.Anything).Return(nil)

			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLastOffset(ctx).Return(wantTail+1, nil)

			require.NoError(t, commands.CreateCheckPoint(ctx))
		})
		t.Run("should return error if failed to restore state", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)

			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			mockAllTimesItems := newMockTopKItems(t)
			topKItemsFactory, _ := mockDeps.TopKItemsFactory.(*mockTopKItemsFactory)
			topKItemsFactory.EXPECT().newTopKItems(topKMaxItemsSize).Return(mockAllTimesItems)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			wantErr := errors.New(faker.Sentence())
			checkPointer.EXPECT().restoreState(ctx, mock.Anything).Return(wantErr)

			require.ErrorIs(t, commands.CreateCheckPoint(ctx), wantErr)
		})
		t.Run("should return error if failed to read lag", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)

			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			mockAllTimesItems := newMockTopKItems(t)
			topKItemsFactory, _ := mockDeps.TopKItemsFactory.(*mockTopKItemsFactory)
			topKItemsFactory.EXPECT().newTopKItems(topKMaxItemsSize).Return(mockAllTimesItems)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, mock.Anything).Return(nil)

			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			wantErr := errors.New(faker.Sentence())
			reader.EXPECT().ReadLastOffset(ctx).Return(0, wantErr)

			require.ErrorIs(t, commands.CreateCheckPoint(ctx), wantErr)
		})
	})
}
