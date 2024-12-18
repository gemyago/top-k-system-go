package aggregation

import (
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/go-faker/faker/v4"
)

func randomTopKItem() *topKItem {
	return &topKItem{
		ItemID: faker.UUIDHyphenated(),
		Count:  1000 + rand.Int64N(100000),
	}
}

func randomTopKItems(n int) []*topKItem {
	items := make([]*topKItem, 0, n)
	for range n {
		items = append(items, randomTopKItem())
	}
	return items
}

func BenchmarkTopKItems(b *testing.B) {
	runTopKItemsTestSuite := func(b *testing.B, newTopKBTreeItems func(maxSize int) topKItems) {
		b.Run("getItems", func(b *testing.B) {
			items := newTopKBTreeItems(1000)
			items.load(randomTopKItems(1000))

			b.Run("get top 100 items", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					items.getItems(100)
				}
			})

			b.Run("get top 500 items", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					items.getItems(500)
				}
			})

			b.Run("get top 1000 items", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					items.getItems(1000)
				}
			})
		})

		b.Run("updateIfGreater", func(b *testing.B) {
			items := newTopKBTreeItems(1000)
			items.load(randomTopKItems(1000))
			maxItem := items.getItems(1000)[0]
			randomItem := items.getItems(1000)[rand.IntN(1000)]

			b.Run("replace existing item", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					item := *randomItem
					item.Count = maxItem.Count + int64(i) + 1
					items.updateIfGreater(item)
				}
			})

			b.Run("replace existing if greater", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					item := *maxItem
					item.ItemID += strconv.Itoa(i)
					item.Count = maxItem.Count + int64(i) + 1
					items.updateIfGreater(item)
				}
			})
		})
	}

	b.Run("topKBTreeItems", func(b *testing.B) {
		runTopKItemsTestSuite(b, func(maxSize int) topKItems {
			return newTopKBTreeItems(maxSize)
		})
	})

	b.Run("topKBTreeItems(sync)", func(b *testing.B) {
		runTopKItemsTestSuite(b, func(maxSize int) topKItems {
			return &synchronisedTopKItems{
				topKItems: newTopKBTreeItems(maxSize),
			}
		})
	})

	b.Run("topKHeapItems", func(b *testing.B) {
		runTopKItemsTestSuite(b, func(maxSize int) topKItems {
			return newTopKHeapItems(maxSize)
		})
	})

	b.Run("topKHeapItems(sync)", func(b *testing.B) {
		runTopKItemsTestSuite(b, func(maxSize int) topKItems {
			return &synchronisedTopKItems{
				topKItems: newTopKHeapItems(maxSize),
			}
		})
	})
}
