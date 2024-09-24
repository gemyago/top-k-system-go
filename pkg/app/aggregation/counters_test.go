package aggregation

import (
	"math/rand"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestCounters(t *testing.T) {
	t.Run("updateItemsCount", func(t *testing.T) {
		t.Run("should update items counters with new values", func(t *testing.T) {
			c := NewCounters()

			initialBaseOffset := rand.Int63n(1000)
			existingData := map[string]int64{
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
			}

			cImpl, _ := c.(*counters)
			cImpl.lastOffset = initialBaseOffset
			for k, v := range existingData {
				cImpl.itemCounters[k] = v
			}

			nextOffset := initialBaseOffset + rand.Int63n(1000)
			newCounts := make(map[string]int64, len(existingData))
			for k := range existingData {
				newCounts[k] = rand.Int63n(1000)
			}

			wantNewData := make(map[string]int64, len(existingData))
			for k, v := range newCounts {
				wantNewData[k] = existingData[k] + v
			}

			c.updateItemsCount(nextOffset, newCounts)
			assert.Equal(t, nextOffset, cImpl.lastOffset)
			assert.Equal(t, wantNewData, cImpl.itemCounters)
		})
	})
}
