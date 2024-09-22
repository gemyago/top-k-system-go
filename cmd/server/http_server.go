package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"time"

	"github.com/gemyago/top-k-system-go/pkg/api/http/routes"
	"github.com/gemyago/top-k-system-go/pkg/api/http/server"
	"github.com/gemyago/top-k-system-go/pkg/di"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

type runParams struct {
	dig.In `ignore-unexported:"true"`

	RootLogger       *slog.Logger
	HTTPServer       *http.Server
	ShutdownHandlers []di.ProcessShutdownHandler `group:"shutdown-handlers"`

	noopHTTPListen bool
}

func run(params runParams) {
	rootLogger := params.RootLogger
	httpServer := params.HTTPServer
	rootCtx := context.Background()

	listenersErrors := make(chan error, 1)
	go func() {
		rootLogger.InfoContext(rootCtx, "Starting http listener",
			slog.String("addr", httpServer.Addr),
			slog.String("idleTimeout", httpServer.IdleTimeout.String()),
			slog.String("readHeaderTimeout", httpServer.ReadHeaderTimeout.String()),
			slog.String("readTimeout", httpServer.ReadTimeout.String()),
			slog.String("writeTimeout", httpServer.WriteTimeout.String()),
		)
		if params.noopHTTPListen {
			rootLogger.InfoContext(rootCtx, "NOOP: Exiting now")
			listenersErrors <- nil
		} else {
			listenersErrors <- httpServer.ListenAndServe()
		}
	}()

	signalCtx, cancel := signal.NotifyContext(rootCtx, unix.SIGINT, unix.SIGTERM)
	defer cancel()

	select {
	case err := <-listenersErrors:
		if err != nil {
			rootLogger.ErrorContext(rootCtx, "Listener error", "err", err)
		} else {
			rootLogger.InfoContext(rootCtx, "Listener stopped")
		}
	case <-signalCtx.Done(): // coverage-ignore
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
}

func newHTTPServerCmd(container *dig.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Command to start http server",
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
		return container.Invoke(func(params runParams) {
			params.noopHTTPListen = noop
			run(params)
		})
	}
	return cmd
}
