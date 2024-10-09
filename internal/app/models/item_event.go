package models

import "time"

type ItemEvent struct {
	ItemID     string    `json:"itemId"`
	IngestedAt time.Time `json:"ingestedAt"`
}
