package aggregation

import (
	"math/rand"
	"testing"

	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/stretchr/testify/assert"
)

func TestAggregatorModel(t *testing.T) {
	newMockDeps := func() ItemEventsAggregatorDeps {
		return ItemEventsAggregatorDeps{}
	}

	t.Run("aggregateItemEvent", func(t *testing.T) {
		t.Run("should set new counters to 1", func(t *testing.T) {
			mockDeps := newMockDeps()
			model := NewItemEventsAggregatorModel(mockDeps)

			baseOffset := rand.Int63()
			itemEvents := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}
			for i, e := range itemEvents {
				model.aggregateItemEvent(baseOffset+int64(i), &e)
			}

			modelImpl, _ := model.(*itemEventsAggregatorModel)
			assert.Equal(t, baseOffset+int64(len(itemEvents)-1), modelImpl.lastAggregatedOffset)
			for _, e := range itemEvents {
				assert.Equal(t, int64(1), modelImpl.aggregatedItems[e.ItemID])
			}
		})
		t.Run("should increment existing counters", func(t *testing.T) {
			mockDeps := newMockDeps()
			model := NewItemEventsAggregatorModel(mockDeps)

			baseCounter := rand.Int63()
			baseOffset := rand.Int63()
			itemEvents := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}
			modelImpl, _ := model.(*itemEventsAggregatorModel)
			for i, e := range itemEvents {
				modelImpl.aggregatedItems[e.ItemID] = baseCounter + int64(i)
				model.aggregateItemEvent(baseOffset+int64(i), &e)
			}

			assert.Equal(t, baseOffset+int64(len(itemEvents)-1), modelImpl.lastAggregatedOffset)
			for i, e := range itemEvents {
				assert.Equal(t, baseCounter+int64(i+1), modelImpl.aggregatedItems[e.ItemID])
			}
		})
	})
}
