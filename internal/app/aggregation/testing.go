//go:build !release

package aggregation

import (
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
)

func randomManifest() checkPointManifest {
	return checkPointManifest{
		LastOffset:           rand.Int64N(10000),
		CountersBlobFileName: faker.Word(),
		AllTimeItemsFileName: faker.Word(),
	}
}

func randomCountersValues() map[string]int64 {
	return map[string]int64{
		faker.UUIDHyphenated(): rand.Int64(),
		faker.UUIDHyphenated(): rand.Int64(),
		faker.UUIDHyphenated(): rand.Int64(),
	}
}
