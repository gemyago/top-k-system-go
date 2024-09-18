package main

import (
	"log/slog"

	"github.com/gemyago/top-k-system-go/pkg/di"
	"github.com/gemyago/top-k-system-go/pkg/diag"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

type rootCmdParams struct {
	container *dig.Container

	childCommands []*cobra.Command
}

func newRootCmd(params rootCmdParams) *cobra.Command {
	verbose := false
	cmd := &cobra.Command{
		Use:   "local-broker",
		Short: "Commands to setup and interact with local broker",
	}
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Produce logs with debug level")
	cmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		logLevel := lo.If(verbose, slog.LevelDebug).Else(slog.LevelInfo)

		rootLogger := diag.SetupRootLogger(
			diag.NewRootLoggerOpts().
				WithJSONLogs(true).
				WithLogLevel(logLevel),
		)

		return params.container.Provide(di.ProvideValue(rootLogger).Constructor)
	}
	for _, child := range params.childCommands {
		cmd.AddCommand(child)
	}
	return cmd
}
