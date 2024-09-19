package ingestion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/segmentio/kafka-go"
	"go.uber.org/dig"
)

type Commands interface {
	IngestItemEvent(ctx context.Context, evt *models.ItemEvent) error
}

type CommandsDeps struct {
	dig.In

	ItemEventsWriter services.ItemEventsKafkaWriter
}

type commands struct {
	CommandsDeps
}

func (c *commands) IngestItemEvent(ctx context.Context, evt *models.ItemEvent) error {
	msgValue, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = c.ItemEventsWriter.WriteMessages(
		// It's going to write in batches outside of the API call
		// we don't want to cancel to abort it
		context.WithoutCancel(ctx),
		kafka.Message{
			Key:   []byte(evt.ItemID),
			Value: msgValue,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write item event (itemID=%v): %w", evt.ItemID, err)
	}
	return nil
}

func NewCommands(deps CommandsDeps) Commands {
	return &commands{deps}
}
