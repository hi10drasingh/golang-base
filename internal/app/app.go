package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/droomlab/drm-coupon/internal/appcontext"
	"github.com/droomlab/drm-coupon/internal/transport/rest"
)

// Run the http server
func Run(configDir string) error {

	var env string = "local"
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()

	AppCtx, err := appcontext.InitilizeAppContext(configDir, env)
	if err != nil {
		return err
	}
	defer AppCtx.Close()

	h := rest.NewHandlers(AppCtx)

	server := h.GetServer()

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	go func() {
		h.AppCtx.Log.Infof("app : API listening on port %v", h.AppCtx.Config.HTTP.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting server: %s", err)

	case <-shutdown:
		h.AppCtx.Log.Info("app : Start shutdown")

		// give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// asking listener to shutdown
		err := server.Shutdown(ctx)
		if err != nil {
			h.AppCtx.Log.Infof("app : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = server.Close()
		}

		if err != nil {
			return fmt.Errorf("app : could not stop server gracefully : %v", err)
		}
	}

	return nil
}
