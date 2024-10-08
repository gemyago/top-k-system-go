package ingestion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gemyago/top-k-system-go/internal/app/models"
	"github.com/segmentio/kafka-go"
	"go.uber.org/dig"
)

type itemEventsWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type CommandsDeps struct {
	dig.In

	ItemEventsWriter itemEventsWriter
}

type Commands struct {
	deps CommandsDeps
}

func (c *Commands) IngestItemEvent(ctx context.Context, evt *models.ItemEvent) error {
	msgValue, err := json.Marshal(evt)
	if err != nil { // coverage-ignore // unrealistic to simulate this error
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = c.deps.ItemEventsWriter.WriteMessages(
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

func NewCommands(deps CommandsDeps) *Commands {
	return &Commands{deps}
}
