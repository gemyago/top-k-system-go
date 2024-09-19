package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gemyago/top-k-system-go/pkg/app/ingestion"
	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

type sendTestEventCmdParams struct {
	container *dig.Container
}

func newSendTestEventCmd(cmdParams sendTestEventCmdParams) *cobra.Command {
	type invokeCmdParams struct {
		dig.In

		RootLogger *slog.Logger

		IngestionCommands ingestion.Commands
	}

	var itemID string
	var eventsNumber int

	cmd := &cobra.Command{
		Use:   "send-test-events",
		Short: "Send test item events",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmdParams.container.Invoke(func(params invokeCmdParams) error {
				logger := params.RootLogger.WithGroup("send-test-event")
				logger.InfoContext(cmd.Context(), "Sending test item event")

				if itemID == "" {
					itemID = lo.Must(uuid.NewV4()).String()
				}
				now := time.Now()
				for range eventsNumber {
					event := models.ItemEvent{
						ItemID:     itemID,
						IngestedAt: now,
						Count:      1,
					}
					if err := params.IngestionCommands.IngestItemEvent(
						cmd.Context(), &event,
					); err != nil {
						return fmt.Errorf("failed to write event: %w", err)
					}
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
	cmd.Flags().StringVar(
		&itemID, "item-id", "", "ItemID to produce the events for. If not provided - random is generated.",
	)
	cmd.Flags().IntVarP(&eventsNumber, "events-number", "n", 1, "Number of events to produce")
	return cmd
}
