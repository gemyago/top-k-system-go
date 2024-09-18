package main

import (
	"log/slog"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

type setupBrokerCmdParams struct {
	container *dig.Container
}

func newSetupBrokerCmd(params setupBrokerCmdParams) *cobra.Command {
	type invokeCmdParams struct {
		dig.In

		RootLogger *slog.Logger
	}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup local broker (topics)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return params.container.Invoke(func(params invokeCmdParams) {
				logger := params.RootLogger.WithGroup("setup-broker")
				logger.InfoContext(cmd.Context(), "Setting up local broker")
			})
		},
	}
	return cmd
}
