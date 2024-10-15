package aggregation

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"math/rand/v2"
	"testing"

	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckPointer(t *testing.T) {
	newMockDeps := func(t *testing.T) CheckPointerDeps {
		return CheckPointerDeps{
			RootLogger:        diag.RootTestLogger(),
			CheckPointerModel: newMockCheckPointerModel(t),
		}
	}

	t.Run("restoreState", func(t *testing.T) {
		t.Run("should read the the manifest and values", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()
			manifest := randomManifest()
			values := randomCountersValues()
			allTimeRawItems := randomTopKItems(10)

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			mockModel.EXPECT().readManifest(ctx).Return(manifest, nil)
			mockModel.EXPECT().readCounters(ctx, manifest.CountersBlobFileName).Return(values, nil)
			mockModel.EXPECT().readItems(ctx, manifest.AllTimeItemsFileName).Return(allTimeRawItems, nil)

			counters, _ := newCounters().(*countersImpl)
			allTimeItems := newTopKItems(topKMaxItemsSize)
			require.NoError(t, cp.restoreState(ctx, aggregationState{
				counters:     counters,
				allTimeItems: allTimeItems,
			}))

			wantAllTimeItems := newTopKItems(topKMaxItemsSize)
			wantAllTimeItems.load(allTimeRawItems)

			assert.Equal(t, manifest.LastOffset, counters.lastOffset)
			assert.Equal(t, values, counters.itemCounters)
			assert.Equal(t,
				wantAllTimeItems.getItems(topKGetAllItemsLimit),
				allTimeItems.getItems(topKGetAllItemsLimit),
			)
		})
		t.Run("should handle initial blank state", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			mockModel.EXPECT().readManifest(ctx).Return(checkPointManifest{}, fmt.Errorf("empty state: %w", fs.ErrNotExist))

			counters, _ := newCounters().(*countersImpl)
			require.NoError(t, cp.restoreState(ctx, aggregationState{
				counters: counters,
			}))

			assert.Equal(t, int64(0), counters.lastOffset)
			assert.Empty(t, counters.itemCounters)
		})
		t.Run("should fail on manifest reading errors", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			wantErr := errors.New(faker.Sentence())
			manifest := randomManifest()

			mockModel.EXPECT().readManifest(ctx).Return(manifest, nil)
			mockModel.EXPECT().readCounters(ctx, manifest.CountersBlobFileName).Return(nil, wantErr)

			counters, _ := newCounters().(*countersImpl)
			require.ErrorIs(t, cp.restoreState(ctx, aggregationState{
				counters: counters,
			}), wantErr)
		})
		t.Run("should fail on counters reading errors", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			wantErr := errors.New(faker.Sentence())
			mockModel.EXPECT().readManifest(ctx).Return(checkPointManifest{}, wantErr)

			counters, _ := newCounters().(*countersImpl)
			require.ErrorIs(t, cp.restoreState(ctx, aggregationState{
				counters: counters,
			}), wantErr)
		})
		t.Run("should fail on all time items reading errors", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()

			manifest := randomManifest()
			values := randomCountersValues()

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			mockModel.EXPECT().readManifest(ctx).Return(manifest, nil)
			mockModel.EXPECT().readCounters(ctx, manifest.CountersBlobFileName).Return(values, nil)
			wantErr := errors.New(faker.Sentence())
			mockModel.EXPECT().readItems(ctx, manifest.AllTimeItemsFileName).Return(nil, wantErr)

			counters, _ := newCounters().(*countersImpl)
			allTimeItems := newTopKItems(topKMaxItemsSize)
			require.ErrorIs(t, cp.restoreState(ctx, aggregationState{
				counters:     counters,
				allTimeItems: allTimeItems,
			}), wantErr)
		})
	})

	t.Run("dumpState", func(t *testing.T) {
		t.Run("should write values and manifest", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()
			values := randomCountersValues()
			cnt := newCounters()
			cnt.updateItemsCount(rand.Int64(), values)
			wantAllTimeItems := newTopKItems(topKMaxItemsSize)
			wantAllTimeItems.load(randomTopKItems(10))

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			mockModel.EXPECT().writeCounters(
				ctx,
				fmt.Sprintf("counters-%d", cnt.getLastOffset()),
				values,
			).Return(nil)
			mockModel.EXPECT().writeItems(
				ctx,
				fmt.Sprintf("all-time-items-%d", cnt.getLastOffset()),
				wantAllTimeItems.getItems(topKMaxItemsSize),
			).Return(nil)
			mockModel.EXPECT().writeManifest(
				ctx,
				checkPointManifest{
					LastOffset:           cnt.getLastOffset(),
					CountersBlobFileName: fmt.Sprintf("counters-%d", cnt.getLastOffset()),
					AllTimeItemsFileName: fmt.Sprintf("all-time-items-%d", cnt.getLastOffset()),
				},
			).Return(nil)

			require.NoError(t, cp.dumpState(ctx, aggregationState{
				counters:     cnt,
				allTimeItems: wantAllTimeItems,
			}))
		})
		t.Run("should handle write counters errors", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()
			values := randomCountersValues()
			cnt := newCounters()
			cnt.updateItemsCount(rand.Int64(), values)

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			wantErr := errors.New(faker.Sentence())
			mockModel.EXPECT().writeCounters(
				ctx,
				fmt.Sprintf("counters-%d", cnt.getLastOffset()),
				values,
			).Return(wantErr)

			require.ErrorIs(t, cp.dumpState(ctx, aggregationState{
				counters: cnt,
			}), wantErr)
		})
		t.Run("should handle write items errors", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()
			values := randomCountersValues()
			cnt := newCounters()
			cnt.updateItemsCount(rand.Int64(), values)
			allTimeItems := newTopKItems(topKMaxItemsSize)
			allTimeItems.load(randomTopKItems(10))

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			mockModel.EXPECT().writeCounters(
				ctx,
				fmt.Sprintf("counters-%d", cnt.getLastOffset()),
				values,
			).Return(nil)
			wantErr := errors.New(faker.Sentence())
			mockModel.EXPECT().writeItems(
				ctx,
				fmt.Sprintf("all-time-items-%d", cnt.getLastOffset()),
				allTimeItems.getItems(topKMaxItemsSize),
			).Return(wantErr)

			require.ErrorIs(t, cp.dumpState(ctx, aggregationState{
				counters:     cnt,
				allTimeItems: allTimeItems,
			}), wantErr)
		})
		t.Run("should handle write manifest errors", func(t *testing.T) {
			deps := newMockDeps(t)
			cp := newCheckPointer(deps)

			ctx := context.Background()
			values := randomCountersValues()
			cnt := newCounters()
			cnt.updateItemsCount(rand.Int64(), values)
			allTimeItems := newTopKItems(topKMaxItemsSize)
			allTimeItems.load(randomTopKItems(10))

			mockModel, _ := deps.CheckPointerModel.(*mockCheckPointerModel)
			wantErr := errors.New(faker.Sentence())
			mockModel.EXPECT().writeCounters(
				ctx,
				fmt.Sprintf("counters-%d", cnt.getLastOffset()),
				values,
			).Return(nil)
			mockModel.EXPECT().writeItems(
				ctx,
				fmt.Sprintf("all-time-items-%d", cnt.getLastOffset()),
				allTimeItems.getItems(topKMaxItemsSize),
			).Return(nil)
			mockModel.EXPECT().writeManifest(
				ctx,
				checkPointManifest{
					LastOffset:           cnt.getLastOffset(),
					CountersBlobFileName: fmt.Sprintf("counters-%d", cnt.getLastOffset()),
					AllTimeItemsFileName: fmt.Sprintf("all-time-items-%d", cnt.getLastOffset()),
				},
			).Return(wantErr)

			require.ErrorIs(t, cp.dumpState(ctx, aggregationState{
				counters:     cnt,
				allTimeItems: allTimeItems,
			}), wantErr)
		})
	})
}
