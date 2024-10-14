package aggregation

import (
	"math/rand/v2"
	"testing"

	"github.com/go-faker/faker/v4"
)

func BenchmarkTopKItems(b *testing.B) {
	randomItem := func() *topKItem {
		return &topKItem{
			itemID: faker.UUIDHyphenated(),
			count:  1000 + rand.Int64N(100000),
		}
	}
	randomItems := func(n int) []*topKItem {
		items := make([]*topKItem, 0, n)
		for range n {
			items = append(items, randomItem())
		}
		return items
	}

	b.Run("getItems", func(b *testing.B) {
		items := newTopKItems(1000)
		items.load(randomItems(1000))

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
	})
}
