package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/signal"
	"time"

	"github.com/gemyago/top-k-system-go/pkg/api/http/routes"
	"github.com/gemyago/top-k-system-go/pkg/api/http/server"
	"github.com/gemyago/top-k-system-go/pkg/app/aggregation"
	"github.com/gemyago/top-k-system-go/pkg/di"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

type createCheckPointParams struct {
	dig.In `ignore-unexported:"true"`

	RootLogger *slog.Logger

	AggregationCommands aggregation.Commands

	ShutdownHandlers []di.ProcessShutdownHandler `group:"shutdown-handlers"`

	noop bool
}

func createCheckPoint(params createCheckPointParams) error {
	rootLogger := params.RootLogger
	rootCtx := context.Background()

	shutdown := func() {
		rootLogger.InfoContext(rootCtx, "Trying to shut down gracefully")
		ts := time.Now()

		grp := errgroup.Group{}
		for _, h := range params.ShutdownHandlers {
			grp.Go(func() error {
				rootLogger.InfoContext(rootCtx, fmt.Sprintf("Shutting down %s", h.Name))
				return h.Shutdown(rootCtx)
			})
		}

		// Not much we can do at this stage, so just logging
		if err := grp.Wait(); err != nil {
			rootLogger.ErrorContext(rootCtx, "Graceful shutdown failed", "err", err)
		}

		rootLogger.InfoContext(rootCtx, "Service stopped",
			slog.Duration("duration", time.Since(ts)),
		)
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
			rootLogger.ErrorContext(rootCtx, "Server error", "err", startupErr)
		} else {
			rootLogger.InfoContext(rootCtx, "Server stopped")
		}
		shutdown()
	case <-signalCtx.Done(): // coverage-ignore
		shutdown()
	}
	return startupErr
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