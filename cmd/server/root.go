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

	// 	port := flag.Int("port", 8080, "Port to listen on")
	// 	jsonLogs := flag.Bool("json-logs", false, "Indicates if logs should be in JSON format or text (default)")
	// 	logLevel := flag.String("log-level", slog.LevelDebug.String(), "Log level can be DEBUG, INFO, WARN and ERROR")
	// 	noop := flag.Bool("noop", false, "Do not start. Just setup deps and exit. Useful for testing if setup is all working.")

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
