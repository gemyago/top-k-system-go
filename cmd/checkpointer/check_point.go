package main

import (
	"context"
	"errors"
	"log/slog"
	"os/signal"
	"time"

	"github.com/gemyago/top-k-system-go/internal/api/http/routes"
	"github.com/gemyago/top-k-system-go/internal/api/http/server"
	"github.com/gemyago/top-k-system-go/internal/app/aggregation"
	"github.com/gemyago/top-k-system-go/internal/di"
	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"golang.org/x/sys/unix"
)

type createCheckPointParams struct {
	dig.In `ignore-unexported:"true"`

	RootLogger *slog.Logger

	AggregationCommands aggregation.Commands

	services.ShutdownHooks

	noop bool
}

func createCheckPoint(params createCheckPointParams) error {
	rootLogger := params.RootLogger
	rootCtx := context.Background()

	shutdown := func() error {
		rootLogger.InfoContext(rootCtx, "Trying to shut down gracefully")
		ts := time.Now()

		err := params.ShutdownHooks.PerformShutdown(rootCtx)
		if err != nil {
			rootLogger.ErrorContext(rootCtx, "Failed to shut down gracefully", diag.ErrAttr(err))
		}

		rootLogger.InfoContext(rootCtx, "Service stopped",
			slog.Duration("duration", time.Since(ts)),
		)
		return err
	}

	signalCtx, cancel := signal.NotifyContext(rootCtx, unix.SIGINT, unix.SIGTERM)
	defer cancel()

	startupErrors := make(chan error)
	go func() {
		if params.noop {
			rootLogger.InfoContext(signalCtx, "NOOP: Exiting now")
			startupErrors <- nil
			return
		}
		startupErrors <- params.AggregationCommands.CreateCheckPoint(signalCtx)
	}()

	var startupErr error
	select {
	case startupErr = <-startupErrors:
		if startupErr != nil {
			rootLogger.ErrorContext(rootCtx, "Server startup failed", "err", startupErr)
		}
	case <-signalCtx.Done(): // coverage-ignore
		// We will attempt to shut down in both cases
		// so doing it once on a next line
	}
	return errors.Join(startupErr, shutdown())
}

func newCreateCheckPointCmd(container *dig.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-check-point",
		Short: "Command to create check point",
	}
	noop := false
	cmd.Flags().BoolVar(
		&noop,
		"noop",
		false,
		"Do not start. Just setup deps and exit. Useful for testing if setup is all working.",
	)
	cmd.PreRunE = func(_ *cobra.Command, _ []string) error {
		return errors.Join(
			// http related dependencies
			routes.Register(container),
			di.ProvideAll(
				container,
				server.NewHTTPServer,
				server.NewRootHandler,
			),
		)
	}
	cmd.RunE = func(_ *cobra.Command, _ []string) error {
		return container.Invoke(func(params createCheckPointParams) error {
			params.noop = noop
			return createCheckPoint(params)
		})
	}
	return cmd
}
