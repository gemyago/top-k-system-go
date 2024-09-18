package models

import "time"

type ItemEvent struct {
	ItemID     string
	IngestedAt time.Time
	Count      int64
}
