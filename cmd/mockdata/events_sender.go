package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gemyago/top-k-system-go/internal/app/models"
	"go.uber.org/dig"
)

type eventsSenderImpl struct {
	// all injectable fields must be exported
	// to let dig inject them

	dig.In

	RootLogger *slog.Logger

	// app layer
	IngestionCommands ingestionCommands
}

func (impl *eventsSenderImpl) sendTestEvent(ctx context.Context, itemID string, eventsNumber int) error {
	evt := &models.ItemEvent{
		ItemID: itemID,
	}
	for range eventsNumber {
		if err := impl.IngestionCommands.IngestItemEvent(ctx, evt); err != nil {
			return fmt.Errorf("failed to ingest item event: %w", err)
		}
	}
	return nil
}

func (impl *eventsSenderImpl) sendTestEvents(ctx context.Context, itemIDsFile string, eventsNumber int) error {
	return nil
}
