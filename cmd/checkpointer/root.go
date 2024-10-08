package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gemyago/top-k-system-go/internal/app/aggregation"
	"github.com/gemyago/top-k-system-go/internal/app/ingestion"
	"github.com/gemyago/top-k-system-go/internal/config"
	"github.com/gemyago/top-k-system-go/internal/di"
	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newRootCmd(container *dig.Container) *cobra.Command {
	logsOutputFile := ""
	jsonLogs := false
	env := ""

	cmd := &cobra.Command{
		Use:   "checkpointer",
		Short: "Command to start checkpointer jobs",
	}
	cmd.PersistentFlags().StringP("log-level", "l", "", "Produce logs with given level. Default is env specific.")
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

		if err = cfg.BindPFlag("defaultLogLevel", cmd.PersistentFlags().Lookup("log-level")); err != nil {
			return err
		}

		var logLevel slog.Level
		if err = logLevel.UnmarshalText([]byte(cfg.GetString("defaultLogLevel"))); err != nil {
			return err
		}

		rootLogger := diag.SetupRootLogger(
			diag.NewRootLoggerOpts().
				WithJSONLogs(jsonLogs).
				WithLogLevel(logLevel).
				WithOptionalOutputFile(logsOutputFile),
		)

		err = errors.Join(
			config.Provide(container, cfg),

			// app layer
			aggregation.Register(container),
			ingestion.Register(container),

			// services
			services.Register(container),

			di.ProvideAll(container,
				di.ProvideValue(rootLogger),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to inject dependencies: %w", err)
		}

		return nil
	}
	return cmd
}
