package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/gemyago/top-k-system-go/pkg/services/blobstorage"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

const defaultItemsNumberToGenerate = 10000

type writeRandomItemsParams struct {
	logger      *slog.Logger
	writer      io.WriteCloser
	itemsNumber int
	services.UUIDGenerator
}

func writeRandomItems(ctx context.Context, params writeRandomItemsParams) error {
	defer params.writer.Close()
	for generated := range params.itemsNumber {
		if generated > 0 && generated%1000000 == 0 {
			params.logger.InfoContext(
				ctx,
				fmt.Sprintf("Generated %d of %d items", generated, params.itemsNumber),
			)
		}
		if _, err := params.writer.Write([]byte(params.UUIDGenerator() + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func newGenerateItemIDsCmd(container *dig.Container) *cobra.Command {
	type invokeCmdParams struct {
		dig.In

		RootLogger *slog.Logger

		// Services
		services.UUIDGenerator
		blobstorage.Storage
	}

	var itemsNumber int
	outputFileName := "test-item-ids.txt"

	noop := false
	cmd := &cobra.Command{
		Use:   "generate-item-ids",
		Short: "Generate test item IDs and write them to file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return container.Invoke(func(params invokeCmdParams) error {
				logger := params.RootLogger.WithGroup("generate-item-ids")
				logger.InfoContext(cmd.Context(), "Generating item ids", slog.Int("number", itemsNumber))

				reader, writer := io.Pipe()

				generatorDone := make(chan error)
				go func() {
					generatorDone <- writeRandomItems(cmd.Context(), writeRandomItemsParams{
						logger:        logger,
						writer:        writer,
						itemsNumber:   itemsNumber,
						UUIDGenerator: params.UUIDGenerator,
					})
				}()

				if err := params.Storage.Upload(cmd.Context(), "test-item-ids.txt", reader); err != nil {
					return err
				}
				if err := <-generatorDone; err != nil {
					return err
				}

				logger.InfoContext(
					cmd.Context(),
					"Test item ids generated",
					slog.Int("number", itemsNumber),
				)

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
	cmd.Flags().StringVarP(&outputFileName, "output-file", "o", outputFileName, "Output file name")
	cmd.Flags().IntVarP(&itemsNumber, "items-number", "n", defaultItemsNumberToGenerate, "Number of items to generate")
	return cmd
}
