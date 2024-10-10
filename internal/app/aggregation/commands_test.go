package aggregation

import (
	"context"
	"errors"
	"math/rand/v2"
	"testing"

	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/go-faker/faker/v4"
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
		}
	}

	t.Run("StartAggregator", func(t *testing.T) {
		t.Run("should restore state and start aggregating", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)
			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*mockItemEventsAggregator)
			aggregator.EXPECT().
				beginAggregating(ctx, wantCounters, beginAggregatingOpts{}).
				Return(nil)

			require.NoError(t, commands.StartAggregator(ctx))
		})

		t.Run("should return error if restore state failed", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)
			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			wantErr := errors.New(faker.Sentence())
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(wantErr)

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

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			wantTail := rand.Int64()
			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLastOffset(ctx).Return(wantTail, nil)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*mockItemEventsAggregator)
			aggregator.EXPECT().
				beginAggregating(ctx, wantCounters, beginAggregatingOpts{
					tillOffset: wantTail - 1,
				}).
				Return(nil)

			checkPointer.EXPECT().dumpState(ctx, wantCounters).Return(nil)

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

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			wantTail := lastOffset + 100
			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLastOffset(ctx).Return(wantTail, nil)
			reader.EXPECT().SetOffset(lastOffset + 1).Return(nil)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*mockItemEventsAggregator)
			aggregator.EXPECT().
				beginAggregating(ctx, wantCounters, beginAggregatingOpts{
					tillOffset: wantTail - 1,
				}).
				Return(nil)

			checkPointer.EXPECT().dumpState(ctx, wantCounters).Return(nil)

			require.NoError(t, commands.CreateCheckPoint(ctx))
		})
		t.Run("should fail if failed to set the offset", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)

			lastOffset := rand.Int64()
			wantCounters.EXPECT().getLastOffset().Return(lastOffset)

			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			wantTail := lastOffset + 100
			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLastOffset(ctx).Return(wantTail, nil)

			wantErr := errors.New(faker.Sentence())
			reader.EXPECT().SetOffset(lastOffset + 1).Return(wantErr)

			require.ErrorIs(t, commands.CreateCheckPoint(ctx), wantErr)
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

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLastOffset(ctx).Return(wantTail+1, nil)
			reader.EXPECT().SetOffset(wantTail + 1).Return(nil)

			require.NoError(t, commands.CreateCheckPoint(ctx))
		})
		t.Run("should return error if failed to restore state", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)

			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			wantErr := errors.New(faker.Sentence())
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(wantErr)

			require.ErrorIs(t, commands.CreateCheckPoint(ctx), wantErr)
		})
		t.Run("should return error if failed to read lag", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := newMockCounters(t)

			countersFactory, _ := mockDeps.CountersFactory.(*mockCountersFactory)
			countersFactory.EXPECT().newCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*mockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			wantErr := errors.New(faker.Sentence())
			reader.EXPECT().ReadLastOffset(ctx).Return(0, wantErr)

			require.ErrorIs(t, commands.CreateCheckPoint(ctx), wantErr)
		})
	})
}
