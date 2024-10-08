package aggregation

import (
	"context"
	"errors"
	"math/rand/v2"
	"testing"

	"github.com/gemyago/top-k-system-go/pkg/diag"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestCommands(t *testing.T) {
	newMockDeps := func(t *testing.T) CommandsDeps {
		return CommandsDeps{
			RootLogger:           diag.RootTestLogger(),
			CheckPointer:         NewMockCheckPointer(t),
			ItemEventsAggregator: NewMockItemEventsAggregator(t),
			ItemEventsReader:     services.NewMockKafkaReader(t),
			CountersFactory:      NewMockCountersFactory(t),
		}
	}

	t.Run("StartAggregator", func(t *testing.T) {
		t.Run("should restore state and start aggregating", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := NewMockCounters(t)
			countersFactory, _ := mockDeps.CountersFactory.(*MockCountersFactory)
			countersFactory.EXPECT().NewCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*MockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*MockItemEventsAggregator)
			aggregator.EXPECT().
				BeginAggregating(ctx, wantCounters, BeginAggregatingOpts{}).
				Return(nil)

			require.NoError(t, commands.StartAggregator(ctx))
			countersFactory.AssertExpectations(t)
			checkPointer.AssertExpectations(t)
			aggregator.AssertExpectations(t)
		})

		t.Run("should return error if restore state failed", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := NewMockCounters(t)
			countersFactory, _ := mockDeps.CountersFactory.(*MockCountersFactory)
			countersFactory.EXPECT().NewCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*MockCheckPointer)
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

			wantCounters := NewMockCounters(t)

			wantCounters.EXPECT().getLastOffset().Return(0)

			countersFactory, _ := mockDeps.CountersFactory.(*MockCountersFactory)
			countersFactory.EXPECT().NewCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*MockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			wantLag := rand.Int64()
			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLag(ctx).Return(wantLag, nil)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*MockItemEventsAggregator)
			aggregator.EXPECT().
				BeginAggregating(ctx, wantCounters, BeginAggregatingOpts{
					TillOffset: wantLag - 1,
				}).
				Return(nil)

			checkPointer.EXPECT().dumpState(ctx, wantCounters).Return(nil)

			require.NoError(t, commands.CreateCheckPoint(ctx))
		})
		t.Run("should set the restored offset of the reader", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := NewMockCounters(t)

			lastOffset := rand.Int64()
			wantCounters.EXPECT().getLastOffset().Return(lastOffset)

			countersFactory, _ := mockDeps.CountersFactory.(*MockCountersFactory)
			countersFactory.EXPECT().NewCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*MockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			wantLag := lastOffset + 100
			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLag(ctx).Return(wantLag, nil)
			reader.EXPECT().SetOffset(lastOffset + 1).Return(nil)

			aggregator, _ := mockDeps.ItemEventsAggregator.(*MockItemEventsAggregator)
			aggregator.EXPECT().
				BeginAggregating(ctx, wantCounters, BeginAggregatingOpts{
					TillOffset: wantLag - 1,
				}).
				Return(nil)

			checkPointer.EXPECT().dumpState(ctx, wantCounters).Return(nil)

			require.NoError(t, commands.CreateCheckPoint(ctx))
		})
		t.Run("should not aggregate if no new messages", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := NewMockCounters(t)

			lastOffset := rand.Int64()
			wantCounters.EXPECT().getLastOffset().Return(lastOffset)

			countersFactory, _ := mockDeps.CountersFactory.(*MockCountersFactory)
			countersFactory.EXPECT().NewCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*MockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			reader.EXPECT().ReadLag(ctx).Return(lastOffset+1, nil)
			reader.EXPECT().SetOffset(lastOffset + 1).Return(nil)

			require.NoError(t, commands.CreateCheckPoint(ctx))
		})
		t.Run("should return error if failed to restore state", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := NewMockCounters(t)

			countersFactory, _ := mockDeps.CountersFactory.(*MockCountersFactory)
			countersFactory.EXPECT().NewCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*MockCheckPointer)
			wantErr := errors.New(faker.Sentence())
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(wantErr)

			require.ErrorIs(t, commands.CreateCheckPoint(ctx), wantErr)
		})
		t.Run("should return error if failed to read lag", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)

			ctx := context.Background()

			wantCounters := NewMockCounters(t)

			countersFactory, _ := mockDeps.CountersFactory.(*MockCountersFactory)
			countersFactory.EXPECT().NewCounters().Return(wantCounters)

			checkPointer, _ := mockDeps.CheckPointer.(*MockCheckPointer)
			checkPointer.EXPECT().restoreState(ctx, wantCounters).Return(nil)

			reader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)
			wantErr := errors.New(faker.Sentence())
			reader.EXPECT().ReadLag(ctx).Return(0, wantErr)

			require.ErrorIs(t, commands.CreateCheckPoint(ctx), wantErr)
		})
	})
}
