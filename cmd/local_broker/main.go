package main

import (
	"fmt"

	"github.com/gemyago/top-k-system-go/pkg/app/ingestion"
	"github.com/gemyago/top-k-system-go/pkg/di"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func mustNoErrors(errs ...error) {
	for i, err := range errs {
		if err != nil {
			panic(fmt.Sprintf("Error %d: %v", i, err))
		}
	}
}

func main() {
	container := dig.New()

	mustNoErrors(
		di.ProvideAll(container,
			// app layer
			ingestion.NewCommands,

			// service layer
			services.NewTimeProvider,
			services.NewItemEventsKafkaWriter,
		),
	)

	rootCmd := newRootCmd(rootCmdParams{
		container: container,
		childCommands: []*cobra.Command{
			newSendTestEventCmd(sendTestEventCmdParams{container}),
		},
	})
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
