package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gemyago/top-k-system-go/pkg/app/ingestion"
	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newSendTestEventCmd(container *dig.Container) *cobra.Command {
	type invokeCmdParams struct {
		dig.In

		RootLogger *slog.Logger

		// app layer
		IngestionCommands ingestion.Commands

		// service layer
		ItemEventsWriter services.ItemEventsKafkaWriter
	}

	var itemID string
	var eventsNumber int

	noop := false
	cmd := &cobra.Command{
		Use:   "send-test-events",
		Short: "Send test item events",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return container.Invoke(func(params invokeCmdParams) error {
				logger := params.RootLogger.WithGroup("send-test-event")
				logger.InfoContext(cmd.Context(), "Sending test item event")

				if itemID == "" {
					itemID = lo.Must(uuid.NewV4()).String()
				}
				now := time.Now()
				for range eventsNumber {
					if noop {
						logger.InfoContext(cmd.Context(), "NOOP: Ingesting event", slog.String("itemID", itemID))
					} else { // coverage-ignore // our test is high level and it's hard cover this step
						event := models.ItemEvent{
							ItemID:     itemID,
							IngestedAt: now,
						}
						if err := params.IngestionCommands.IngestItemEvent(
							cmd.Context(), &event,
						); err != nil {
							return fmt.Errorf("failed to write event: %w", err)
						}
					}
				}

				if err := params.ItemEventsWriter.Close(); err != nil {
					return fmt.Errorf("failed to flush pending events: %w", err)
				}

				logger.InfoContext(
					cmd.Context(),
					"Test events sent",
					slog.Int("number", eventsNumber),
					slog.String("itemId", itemID),
				)

				return nil
			})
		},
	}
	cmd.Flags().BoolVar(
		&noop,
		"noop",
		false,
		"Do not send. Just setup deps and exit. Useful for testing if setup is all working.",
	)
	cmd.Flags().StringVar(
		&itemID, "item-id", "", "ItemID to produce the events for. If not provided - random is generated.",
	)
	cmd.Flags().IntVarP(&eventsNumber, "events-number", "n", 1, "Number of events to produce")
	return cmd
}
