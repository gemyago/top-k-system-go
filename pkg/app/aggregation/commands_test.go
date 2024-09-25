package aggregation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommands(t *testing.T) {
	newMockDeps := func(t *testing.T) CommandsDeps {
		return CommandsDeps{
			CheckPointer:         NewMockCheckPointer(t),
			ItemEventsAggregator: NewMockItemEventsAggregator(t),
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
	})
}
