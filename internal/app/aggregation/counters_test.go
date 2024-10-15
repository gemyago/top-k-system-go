package aggregation

import (
	"maps"
	"math/rand"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestCounters(t *testing.T) {
	t.Run("updateItemsCount", func(t *testing.T) {
		t.Run("should set items counters with new values", func(t *testing.T) {
			// Doing factory func just to test it
			c := countersFactoryFunc(newCounters).newCounters()

			initialBaseOffset := rand.Int63n(1000)
			newCounts := map[string]int64{
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
			}

			cImpl, _ := c.(*countersImpl)
			cImpl.lastOffset = initialBaseOffset

			nextOffset := initialBaseOffset + rand.Int63n(1000)
			gotUpdated := c.updateItemsCount(nextOffset, newCounts)
			assert.Equal(t, nextOffset, cImpl.lastOffset)
			assert.Equal(t, newCounts, cImpl.itemCounters)
			assert.Equal(t, newCounts, gotUpdated)
		})

		t.Run("should increment existing items counters", func(t *testing.T) {
			// Doing factory func just to test it
			c := countersFactoryFunc(newCounters).newCounters()

			initialBaseOffset := rand.Int63n(1000)
			existingNonUpdatableData := map[string]int64{
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
			}
			existingUpdatableData := map[string]int64{
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
				faker.UUIDHyphenated(): rand.Int63n(1000),
			}
			existingData := make(map[string]int64, len(existingNonUpdatableData)+len(existingUpdatableData))
			maps.Copy(existingData, existingNonUpdatableData)
			maps.Copy(existingData, existingUpdatableData)

			cImpl, _ := c.(*countersImpl)
			cImpl.lastOffset = initialBaseOffset
			for k, v := range existingData {
				cImpl.itemCounters[k] = v
			}

			nextOffset := initialBaseOffset + rand.Int63n(1000)
			newCounts := make(map[string]int64, len(existingUpdatableData))
			for k := range existingUpdatableData {
				newCounts[k] = rand.Int63n(1000)
			}

			wantResult := make(map[string]int64, len(existingUpdatableData))
			for k, v := range newCounts {
				wantResult[k] = existingData[k] + v
			}

			wantUpdatedData := make(map[string]int64, len(existingData))
			maps.Copy(wantUpdatedData, existingNonUpdatableData)
			maps.Copy(wantUpdatedData, wantResult)

			gotUpdated := c.updateItemsCount(nextOffset, newCounts)
			assert.Equal(t, nextOffset, cImpl.lastOffset)
			assert.Equal(t, wantUpdatedData, cImpl.itemCounters)
			assert.Equal(t, wantResult, gotUpdated)
		})
	})
}
