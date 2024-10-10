package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gemyago/top-k-system-go/internal/app/models"
	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

type ingestionCommands interface {
	IngestItemEvent(ctx context.Context, evt *models.ItemEvent) error
}

func newSendTestEventCmd(container *dig.Container) *cobra.Command {
	type invokeCmdParams struct {
		dig.In

		RootLogger *slog.Logger

		// app layer
		IngestionCommands ingestionCommands

		// service layer
		ItemEventsWriter services.ItemEventsKafkaWriter

		// package internal
		EventsSender eventsSender
	}

	var itemID string
	var itemIDsFile string
	var eventsNumber int
	var eventsNumberMax int
	const eventsNumberMaxDefault = 10
	noop := false

	doSend := func(ctx context.Context, params invokeCmdParams) error {
		if itemIDsFile != "" {
			if err := params.EventsSender.sendTestEvents(
				ctx,
				itemIDsFile,
				eventsNumber,
				lo.If(eventsNumberMax == 0, eventsNumber+eventsNumberMaxDefault).Else(eventsNumberMax),
			); err != nil {
				return fmt.Errorf("failed to send test event: %w", err)
			}
		} else {
			if err := params.EventsSender.sendTestEvent(ctx, itemID, eventsNumber); err != nil {
				return fmt.Errorf("failed to send test event: %w", err)
			}
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:   "send-test-events",
		Short: "Send test item events",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := container.Decorate(func(rootLogger *slog.Logger, sender eventsSender) eventsSender {
				return newNoopEventsSender(rootLogger, sender, noop)
			}); err != nil {
				return fmt.Errorf("failed to decorate events sender: %w", err)
			}

			return container.Invoke(func(params invokeCmdParams) error {
				logger := params.RootLogger.WithGroup("send-test-event")
				logger.InfoContext(cmd.Context(), "Sending test item event")

				if itemID == "" {
					itemID = lo.Must(uuid.NewV4()).String()
				}
				if noop {
					logger.InfoContext(
						cmd.Context(),
						"NOOP: Producing test events",
						slog.String("itemID", itemID),
						slog.String("itemIDsFile", itemIDsFile),
					)
				} else {
					return doSend(cmd.Context(), params)
				}

				if err := params.ItemEventsWriter.Close(); err != nil {
					return fmt.Errorf("failed to flush pending events: %w", err)
				}

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
	cmd.Flags().StringVar(
		&itemIDsFile,
		"item-ids-file",
		"",
		"File name with generated item IDs to produce events for (alternative to item-id).",
	)
	cmd.Flags().IntVarP(&eventsNumber, "events-number", "n", 1, "Number of events to produce")
	cmd.Flags().IntVarP(&eventsNumberMax,
		"events-number-max",
		"m", 0,
		"If provided, will generate random number of events between n and m (for file mode only)."+
			" If not provided - n + 10.",
	)
	return cmd
}
