package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gemyago/top-k-system-go/internal/app/models"
	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/gemyago/top-k-system-go/internal/services/blobstorage"
	"go.uber.org/dig"
)

type randIntN func(n int) int

type eventsSender interface {
	sendTestEvent(ctx context.Context, itemID string, eventsNumber int) error
	sendTestEvents(ctx context.Context,
		itemIDsFile string,
		eventsMin int,
		eventsMax int,
	) error
}

type defaultEventsSender struct {
	// all injectable fields must be exported
	// to let dig inject them

	dig.In

	RootLogger *slog.Logger

	// app layer
	IngestionCommands ingestionCommands

	// service layer
	Time services.TimeProvider
	blobstorage.Storage

	// package internal
	RandIntN randIntN
}

func (impl *defaultEventsSender) sendTestEvent(ctx context.Context, itemID string, eventsNumber int) error {
	evt := &models.ItemEvent{
		ItemID:     itemID,
		IngestedAt: impl.Time.Now(),
	}
	for range eventsNumber {
		if err := impl.IngestionCommands.IngestItemEvent(ctx, evt); err != nil {
			return fmt.Errorf("failed to ingest item event: %w", err)
		}
	}
	impl.RootLogger.DebugContext(
		ctx,
		"Test events sent",
		slog.Int("number", eventsNumber),
		slog.String("itemId", itemID),
	)
	return nil
}

func (impl *defaultEventsSender) sendTestEvents(
	ctx context.Context,
	itemIDsFile string,
	eventsMin int,
	eventsMax int,
) error {
	var data bytes.Buffer
	if err := impl.Storage.Download(ctx, itemIDsFile, &data); err != nil {
		return fmt.Errorf("failed to download item IDs from file %s: %w", itemIDsFile, err)
	}
	itemIDs := strings.Split(strings.Trim(data.String(), ""), "\n")
	for _, itemID := range itemIDs {
		eventsNumber := impl.RandIntN(eventsMax-eventsMin) + eventsMin
		if err := impl.sendTestEvent(ctx, itemID, eventsNumber); err != nil {
			return fmt.Errorf("failed to send test events for item %s: %w", itemID, err)
		}
	}
	impl.RootLogger.InfoContext(
		ctx,
		"Test events sent",
		slog.Int("itemIDsNumber", len(itemIDs)),
		slog.Int("eventsMin", eventsMin),
		slog.Int("eventsMax", eventsMax),
	)
	return nil
}

var _ eventsSender = &defaultEventsSender{}

type noopEventsSender struct {
	noop   bool
	target eventsSender
	logger *slog.Logger
}

func (s *noopEventsSender) sendTestEvent(ctx context.Context, itemID string, eventsNumber int) error {
	if s.noop {
		s.logger.InfoContext(ctx, "NOOP: sending test event",
			slog.String("itemID", itemID),
			slog.Int("eventsNumber", eventsNumber),
		)
		return nil
	}
	return s.target.sendTestEvent(ctx, itemID, eventsNumber)
}

func (s *noopEventsSender) sendTestEvents(
	ctx context.Context,
	itemIDsFile string,
	eventsMin int,
	eventsMax int,
) error {
	if s.noop {
		s.logger.InfoContext(ctx, "NOOP: sending test events",
			slog.String("itemIDsFile", itemIDsFile),
			slog.Int("eventsMin", eventsMin),
			slog.Int("eventsMax", eventsMax),
		)
		return nil
	}
	return s.target.sendTestEvents(ctx, itemIDsFile, eventsMin, eventsMax)
}

func newNoopEventsSender(rootLogger *slog.Logger, target eventsSender, noop bool) eventsSender {
	return &noopEventsSender{
		target: target,
		noop:   noop,
		logger: rootLogger.WithGroup("noop-events-sender"),
	}
}
