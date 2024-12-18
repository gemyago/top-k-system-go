package aggregation

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTopKItems(t *testing.T) {
	topKItemsTestSuite := func(t *testing.T, newTopKBTreeItems func(maxSize int) topKItems) {
		t.Run("load", func(t *testing.T) {
			t.Run("should load given items", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}

				items := newTopKBTreeItems(100)
				items.load(originalItems)

				actualItems := items.getItems(100)
				require.Len(t, actualItems, len(originalItems))

				for i, item := range originalItems {
					wantItem := actualItems[len(originalItems)-i-1]
					assert.Equal(t, wantItem.ItemID, item.ItemID)
					assert.Equal(t, wantItem.Count, item.Count)
				}
			})

			t.Run("should load given items and keep only top k items", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems) / 2

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				actualItems := items.getItems(len(originalItems))
				require.Len(t, actualItems, wantItemsCount)

				for i, item := range originalItems[len(originalItems)-wantItemsCount:] {
					wantItem := actualItems[len(actualItems)-i-1]
					assert.Equal(t, wantItem.ItemID, item.ItemID)
					assert.Equal(t, wantItem.Count, item.Count)
				}
			})
		})

		t.Run("getItems", func(t *testing.T) {
			t.Run("should return all items in descending order", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}

				items := newTopKBTreeItems(100)
				items.load(originalItems)

				actualItems := items.getItems(100)
				wantItems := slices.Clone(originalItems)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.Count - i.Count)
				})
				assert.Equal(t, wantItems, actualItems)
			})
			t.Run("should include items with duplicate keys", func(t *testing.T) {
				var baseCount int64 = 10000
				item23Count := baseCount + 100 + rand.Int64N(10)
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: item23Count},
					{ItemID: "item3-" + faker.Word(), Count: item23Count},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}

				items := newTopKBTreeItems(100)
				items.load(originalItems)

				actualItems := items.getItems(100)
				wantItems := slices.Clone(originalItems)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					if i.Count == j.Count {
						return strings.Compare(j.ItemID, i.ItemID)
					}
					return int(j.Count - i.Count)
				})
				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should return all items items using constant", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}

				items := newTopKBTreeItems(100)
				items.load(originalItems)

				actualItems := items.getItems(topKGetAllItemsLimit)
				wantItems := slices.Clone(originalItems)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.Count - i.Count)
				})
				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should return limited items list", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems) / 2

				items := newTopKBTreeItems(100)
				items.load(originalItems)

				actualItems := items.getItems(wantItemsCount)
				wantItems := slices.Clone(originalItems)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.Count - i.Count)
				})
				wantItems = wantItems[:wantItemsCount]
				assert.Equal(t, wantItems, actualItems)
			})
		})

		t.Run("updateIfGreater", func(t *testing.T) {
			t.Run("should insert new item", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems) + 1

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				newItem := topKItem{ItemID: "item6-" + faker.Word(), Count: baseCount + 500 + rand.Int64N(10)}
				items.updateIfGreater(newItem)

				wantItems := slices.Clone(originalItems)
				wantItems = append(wantItems, &newItem)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.Count - i.Count)
				})

				actualItems := items.getItems(len(originalItems) + 1)

				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should update existing item", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems)

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				item3 := *originalItems[2]
				item3.Count *= 100 // should become biggest
				items.updateIfGreater(item3)

				wantItems := slices.Clone(originalItems)
				wantItems[2] = &item3
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.Count - i.Count)
				})

				actualItems := items.getItems(len(originalItems))

				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should update existing item with same count", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems)

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				item3 := *originalItems[2]
				item3.Count = originalItems[1].Count
				items.updateIfGreater(item3)

				wantItems := slices.Clone(originalItems)
				wantItems[2] = &item3
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					if i.Count == j.Count {
						return strings.Compare(j.ItemID, i.ItemID)
					}
					return int(j.Count - i.Count)
				})

				actualItems := items.getItems(len(originalItems))
				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should ignore new item if it's not in top k", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems)

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				newItem := topKItem{ItemID: "item6-" + faker.Word(), Count: baseCount - rand.Int64N(10)}
				items.updateIfGreater(newItem)

				wantItems := slices.Clone(originalItems)
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.Count - i.Count)
				})

				actualItems := items.getItems(len(originalItems))

				assert.Equal(t, wantItems, actualItems)
			})

			t.Run("should replace existing item if new item is in top k", func(t *testing.T) {
				var baseCount int64 = 10000
				originalItems := []*topKItem{
					{ItemID: "item1-" + faker.Word(), Count: baseCount + rand.Int64N(10)},
					{ItemID: "item2-" + faker.Word(), Count: baseCount + 100 + rand.Int64N(10)},
					{ItemID: "item3-" + faker.Word(), Count: baseCount + 200 + rand.Int64N(10)},
					{ItemID: "item4-" + faker.Word(), Count: baseCount + 300 + rand.Int64N(10)},
					{ItemID: "item5-" + faker.Word(), Count: baseCount + 400 + rand.Int64N(10)},
				}
				wantItemsCount := len(originalItems)

				items := newTopKBTreeItems(wantItemsCount)
				items.load(originalItems)

				newItem := topKItem{ItemID: "item6-" + faker.Word(), Count: baseCount + 50 + rand.Int64N(10)}
				items.updateIfGreater(newItem)

				wantItems := slices.Clone(originalItems)
				wantItems[0] = &newItem
				slices.SortFunc(wantItems, func(i, j *topKItem) int {
					return int(j.Count - i.Count)
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

	t.Run("topKHeapItems", func(t *testing.T) {
		topKItemsTestSuite(t, func(maxSize int) topKItems {
			return newTopKHeapItems(maxSize)
		})
	})

	t.Run("synchronisedTopKItems", func(t *testing.T) {
		topKItemsTestSuite(t, func(maxSize int) topKItems {
			return &synchronisedTopKItems{
				topKItems: newTopKBTreeItems(maxSize),
			}
		})
	})
}
