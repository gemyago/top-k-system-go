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
	logsOutputFile := ""
	jsonLogs := false
	env := ""

	cmd := &cobra.Command{
		Use:   "local-broker",
		Short: "Commands to setup and interact with local broker",
	}
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Produce logs with debug level")
	cmd.PersistentFlags().StringVar(
		&logsOutputFile,
		"logs-file",
		"",
		"Produce logs to file instead of stdout. Used for tests only.",
	)
	cmd.PersistentFlags().BoolVar(
		&jsonLogs,
		"json-logs",
		false,
		"Indicates if logs should be in JSON format or text (default)",
	)
	cmd.PersistentFlags().StringVarP(
		&env,
		"env",
		"e",
		"",
		"Env that the process is running in.",
	)
	cmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		cfg, err := config.Load(config.NewLoadOpts().WithEnv(env))
		if err != nil {
			return err
		}

		logLevel := lo.If(verbose, slog.LevelDebug).Else(slog.LevelInfo)

		rootLogger := diag.SetupRootLogger(
			diag.NewRootLoggerOpts().
				WithJSONLogs(jsonLogs).
				WithLogLevel(logLevel).
				WithOptionalOutputFile(logsOutputFile),
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
