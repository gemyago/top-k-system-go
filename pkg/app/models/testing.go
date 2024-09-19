//go:build !release

package models

import (
	"time"

	"github.com/go-faker/faker/v4"
)

func MakeRandomItemEvent() ItemEvent {
	return ItemEvent{
		ItemID:     faker.UUIDHyphenated(),
		IngestedAt: time.UnixMilli(faker.RandomUnixTime()),
	}
}
