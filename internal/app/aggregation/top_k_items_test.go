package aggregation

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestTopKItems(t *testing.T) {
	topKItemsTestSuite := func(t *testing.T, newTopKBTreeItems func(maxSize int) topKItems) {
		t.Run("load", func(t *testing.T) {
			t.Run("should load given items", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{itemID: "item1-" + faker.Word(), count: baseCount + rand.Int64N(10)},
					{itemID: "item2-" + faker.Word(), count: baseCount + 100 + rand.Int64N(10)},
					{itemID: "item3-" + faker.Word(), count: baseCount + 200 + rand.Int64N(10)},
					{itemID: "item4-" + faker.Word(), count: baseCount + 300 + rand.Int64N(10)},
					{itemID: "item5-" + faker.Word(), count: baseCount + 400 + rand.Int64N(10)},
				}

				items := newTopKBTreeItems(100)
				items.load(originalItems)

				actualItems := items.getItems(100)
				assert.Len(t, actualItems, len(originalItems))

				for i, item := range originalItems {
					wantItem := actualItems[len(originalItems)-i-1]
					assert.Equal(t, wantItem.itemID, item.itemID)
					assert.Equal(t, wantItem.count, item.count)
				}
			})

			t.Run("should load given items and keep only top k items", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{itemID: "item1-" + faker.Word(), count: baseCount + rand.Int64N(10)},
					{itemID: "item2-" + faker.Word(), count: baseCount + 100 + rand.Int64N(10)},
					{itemID: "item3-" + faker.Word(), count: baseCount + 200 + rand.Int64N(10)},
					{itemID: "item4-" + faker.Word(), count: baseCount + 300 + rand.Int64N(10)},
					{itemID: "item5-" + faker.Word(), count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems) / 2

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				actualItems := items.getItems(len(originalItems))
				assert.Len(t, actualItems, wantItemsCount)

				for i, item := range originalItems[len(originalItems)-wantItemsCount:] {
					wantItem := actualItems[len(actualItems)-i-1]
					assert.Equal(t, wantItem.itemID, item.itemID)
					assert.Equal(t, wantItem.count, item.count)
				}
			})
		})

		t.Run("getItems", func(t *testing.T) {
			t.Run("should return all items in descending order", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{itemID: "item1-" + faker.Word(), count: baseCount + rand.Int64N(10)},
					{itemID: "item2-" + faker.Word(), count: baseCount + 100 + rand.Int64N(10)},
					{itemID: "item3-" + faker.Word(), count: baseCount + 200 + rand.Int64N(10)},
					{itemID: "item4-" + faker.Word(), count: baseCount + 300 + rand.Int64N(10)},
					{itemID: "item5-" + faker.Word(), count: baseCount + 400 + rand.Int64N(10)},
				}

				items := newTopKBTreeItems(100)
				items.load(originalItems)

				actualItems := items.getItems(100)
				wantItems := slices.Clone(originalItems)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.count - i.count)
				})
				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should return limited items list", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{itemID: "item1-" + faker.Word(), count: baseCount + rand.Int64N(10)},
					{itemID: "item2-" + faker.Word(), count: baseCount + 100 + rand.Int64N(10)},
					{itemID: "item3-" + faker.Word(), count: baseCount + 200 + rand.Int64N(10)},
					{itemID: "item4-" + faker.Word(), count: baseCount + 300 + rand.Int64N(10)},
					{itemID: "item5-" + faker.Word(), count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems) / 2

				items := newTopKBTreeItems(100)
				items.load(originalItems)

				actualItems := items.getItems(wantItemsCount)
				wantItems := slices.Clone(originalItems)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.count - i.count)
				})
				wantItems = wantItems[:wantItemsCount]
				assert.Equal(t, wantItems, actualItems)
			})
		})

		t.Run("updateIfGreater", func(t *testing.T) {
			t.Run("should insert new item", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{itemID: "item1-" + faker.Word(), count: baseCount + rand.Int64N(10)},
					{itemID: "item2-" + faker.Word(), count: baseCount + 100 + rand.Int64N(10)},
					{itemID: "item3-" + faker.Word(), count: baseCount + 200 + rand.Int64N(10)},
					{itemID: "item4-" + faker.Word(), count: baseCount + 300 + rand.Int64N(10)},
					{itemID: "item5-" + faker.Word(), count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems) + 1

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				newItem := topKItem{itemID: "item6-" + faker.Word(), count: baseCount + 500 + rand.Int64N(10)}
				items.updateIfGreater(newItem)

				wantItems := slices.Clone(originalItems)
				wantItems = append(wantItems, &newItem)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.count - i.count)
				})

				actualItems := items.getItems(len(originalItems) + 1)

				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should update existing item", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{itemID: "item1-" + faker.Word(), count: baseCount + rand.Int64N(10)},
					{itemID: "item2-" + faker.Word(), count: baseCount + 100 + rand.Int64N(10)},
					{itemID: "item3-" + faker.Word(), count: baseCount + 200 + rand.Int64N(10)},
					{itemID: "item4-" + faker.Word(), count: baseCount + 300 + rand.Int64N(10)},
					{itemID: "item5-" + faker.Word(), count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems)

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				item3 := *originalItems[2]
				item3.count *= 100 // should become biggest
				items.updateIfGreater(item3)

				wantItems := slices.Clone(originalItems)
				wantItems[2] = &item3
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.count - i.count)
				})

				actualItems := items.getItems(len(originalItems))

				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should ignore new item if it's not in top k", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{itemID: "item1-" + faker.Word(), count: baseCount + rand.Int64N(10)},
					{itemID: "item2-" + faker.Word(), count: baseCount + 100 + rand.Int64N(10)},
					{itemID: "item3-" + faker.Word(), count: baseCount + 200 + rand.Int64N(10)},
					{itemID: "item4-" + faker.Word(), count: baseCount + 300 + rand.Int64N(10)},
					{itemID: "item5-" + faker.Word(), count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems)

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				newItem := topKItem{itemID: "item6-" + faker.Word(), count: baseCount - rand.Int64N(10)}
				items.updateIfGreater(newItem)

				wantItems := slices.Clone(originalItems)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.count - i.count)
				})

				actualItems := items.getItems(len(originalItems))

				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should replace existing item if new item is in top k", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{itemID: "item1-" + faker.Word(), count: baseCount + rand.Int64N(10)},
					{itemID: "item2-" + faker.Word(), count: baseCount + 100 + rand.Int64N(10)},
					{itemID: "item3-" + faker.Word(), count: baseCount + 200 + rand.Int64N(10)},
					{itemID: "item4-" + faker.Word(), count: baseCount + 300 + rand.Int64N(10)},
					{itemID: "item5-" + faker.Word(), count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems)

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				newItem := topKItem{itemID: "item6-" + faker.Word(), count: baseCount + 50 + rand.Int64N(10)}
				items.updateIfGreater(newItem)

				wantItems := slices.Clone(originalItems)
				wantItems[0] = &newItem
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.count - i.count)
				})

				actualItems := items.getItems(len(originalItems))

				assert.Equal(t, wantItems, actualItems)
			})
		})
	}

	t.Run("topKBTreeItems", func(t *testing.T) {
		topKItemsTestSuite(t, func(maxSize int) topKItems {
			return newTopKBTreeItems(maxSize)
		})
	})
}
