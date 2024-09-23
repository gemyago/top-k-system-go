package aggregation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemEventsAggregator(t *testing.T) {
	t.Run("BeginAggregating", func(t *testing.T) {
		t.Run("should exit when context cancelled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			aggregator := NewItemEventsAggregator()

			exit := make(chan error)
			go func() {
				exit <- aggregator.BeginAggregating(ctx)
			}()
			cancel()
			gotErr := <-exit
			assert.NoError(t, gotErr)
		})
	})
}
