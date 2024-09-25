package aggregation

import (
	"context"
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
			vals := randomCountersValues()

			mockModel, _ := deps.CheckPointerModel.(*MockCheckPointerModel)
			mockModel.EXPECT().readManifest(ctx).Return(manifest, nil)
			mockModel.EXPECT().readCounters(ctx, manifest.CountersBlobFileName).Return(vals, nil)

			counters, _ := NewCounters().(*counters)
			require.NoError(t, checkPointer.restoreState(ctx, counters))

			assert.Equal(t, manifest.LastOffset, counters.lastOffset)
			assert.Equal(t, vals, counters.itemCounters)
		})
	})
}
