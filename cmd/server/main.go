package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func setupCommands() *cobra.Command {
	container := dig.New()
	rootCmd := newRootCmd(container)
	rootCmd.AddCommand(
		newHTTPServerCmd(container),
	)
	return rootCmd
}

func main() { // coverage-ignore
	rootCmd := setupCommands()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
