package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/droomlab/drm-coupon/internal/app/handlers"
	"github.com/droomlab/drm-coupon/internal/app/server"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/pkg/errors"
)

func main() {
	logger := drmlog.NewConsoleLogger()
	if err := run(); err != nil {
		logger.Fatal(context.Background(), err, "Main Error")
	}
}

// Run the http server.
func run() (err error) {
	dependencies, err := server.Init()
	if err != nil {
		return errors.Wrap(err, "Initializing dependencies")
	}

	defer func() {
		if er := dependencies.Close(); er != nil {
			err = er
		}
	}()

	// channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	app := handlers.NewHandlers(handlers.Config{
		Shutdown: shutdown,
		Deps:     dependencies,
	})

	srv := server.New(app, dependencies)

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	go func() {
		dependencies.Log.Infof(context.Background(), "app : API listening on port %v", dependencies.Config.HTTP.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// blocking run and waiting for shutdown.
	select {
	case err = <-serverErrors:
		return fmt.Errorf("error: starting server: %w", err)

	case sig := <-shutdown:
		dependencies.Log.Infof(context.Background(), "app : Start shutdown | signal : %v", sig)

		// give outstanding requests a deadline for completion.
		timeout := dependencies.Config.HTTP.ShutdownTimeout.Time
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		defer cancel()

		// asking listener to shutdown
		err = srv.Shutdown(ctx)
		if err != nil {
			dependencies.Log.Infof(ctx, "app : Graceful shutdown did not complete in %v : %w", timeout, err)
			err = srv.Close()
		}

		if err != nil {
			return fmt.Errorf("app : could not stop server gracefully : %w", err)
		}
	}

	return errors.Wrap(err, "Main Run")
}
