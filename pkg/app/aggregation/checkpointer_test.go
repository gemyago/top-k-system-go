package aggregation

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckPointer(t *testing.T) {
	newMockDeps := func(t *testing.T) CheckPointerDeps {
		return CheckPointerDeps{
			CheckPointerModel: NewMockCheckPointerModel(t),
		}
	}

	t.Run("restoreState", func(t *testing.T) {
		t.Run("should read the the manifest and values", func(t *testing.T) {
			deps := newMockDeps(t)
			checkPointer := NewCheckPointer(deps)

			ctx := context.Background()
			manifest := randomManifest()
			values := randomCountersValues()

			mockModel, _ := deps.CheckPointerModel.(*MockCheckPointerModel)
			mockModel.EXPECT().readManifest(ctx).Return(manifest, nil)
			mockModel.EXPECT().readCounters(ctx, manifest.CountersBlobFileName).Return(values, nil)

			counters, _ := NewCounters().(*counters)
			require.NoError(t, checkPointer.restoreState(ctx, counters))

			assert.Equal(t, manifest.LastOffset, counters.lastOffset)
			assert.Equal(t, values, counters.itemCounters)
		})
	})

	t.Run("dumpState", func(t *testing.T) {
		t.Run("should write values and manifest", func(t *testing.T) {
			deps := newMockDeps(t)
			checkPointer := NewCheckPointer(deps)

			ctx := context.Background()
			values := randomCountersValues()
			counters := NewCounters()
			counters.updateItemsCount(rand.Int64(), values)

			mockModel, _ := deps.CheckPointerModel.(*MockCheckPointerModel)
			mockModel.EXPECT().writeCounters(
				ctx,
				fmt.Sprintf("counters-%d", counters.getLastOffset()),
				values,
			).Return(nil)
			mockModel.EXPECT().writeManifest(
				ctx,
				checkPointManifest{
					LastOffset:           counters.getLastOffset(),
					CountersBlobFileName: fmt.Sprintf("counters-%d", counters.getLastOffset()),
				},
			).Return(nil)

			require.NoError(t, checkPointer.dumpState(ctx, counters))
		})
	})
}
