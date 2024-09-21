package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gemyago/top-k-system-go/config"
	"github.com/gemyago/top-k-system-go/pkg/app/ingestion"
	"github.com/gemyago/top-k-system-go/pkg/di"
	"github.com/gemyago/top-k-system-go/pkg/diag"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newRootCmd(container *dig.Container) *cobra.Command {
	verbose := false

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Command to start the server",
	}
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Produce logs with debug level")

	cmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		logLevel := lo.If(verbose, slog.LevelDebug).Else(slog.LevelInfo)

		rootLogger := diag.SetupRootLogger(
			diag.NewRootLoggerOpts().
				WithJSONLogs(true).
				WithLogLevel(logLevel),
		)

		err = errors.Join(
			config.Provide(container, cfg),
			di.ProvideAll(container,
				di.ProvideValue(rootLogger),

				// app layer
				ingestion.NewCommands,

				// service layer
				services.NewTimeProvider,
				services.NewItemEventsKafkaWriter,
			),
		)
		if err != nil {
			return fmt.Errorf("failed to inject dependencies: %w", err)
		}

		return nil
	}
	return cmd
}
