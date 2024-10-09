package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/gemyago/top-k-system-go/internal/services/blobstorage"
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
		if generated > 0 && generated%1000000 == 0 { // coverage-ignore // no value to test log message
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

	itemsNumber := defaultItemsNumberToGenerate
	outputFileName := "test-item-ids.txt"
	overwrite := false

	noop := false
	cmd := &cobra.Command{
		Use:   "generate-item-ids",
		Short: "Generate test item IDs and write them to file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return container.Invoke(func(params invokeCmdParams) error {
				logger := params.RootLogger.WithGroup("generate-item-ids")
				logger.InfoContext(cmd.Context(), "Generating item ids", slog.Int("number", itemsNumber))

				if overwrite {
					logger.InfoContext(cmd.Context(), "Removing existing file", slog.String("file", outputFileName))
					if err := params.Storage.Delete(cmd.Context(), outputFileName); err != nil {
						if !errors.Is(err, os.ErrNotExist) {
							return err
						}
					}
				}

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

				if err := params.Storage.Upload(cmd.Context(), outputFileName, reader); err != nil {
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
	cmd.Flags().IntVarP(&itemsNumber, "items-number", "n", itemsNumber, "Number of items to generate")
	cmd.Flags().BoolVar(&overwrite, "overwrite", overwrite, "Remove existing file before writing")
	return cmd
}
