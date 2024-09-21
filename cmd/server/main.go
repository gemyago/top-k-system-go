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

func executeRootCommand(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func main() { // coverage-ignore
	executeRootCommand(setupCommands())
}
