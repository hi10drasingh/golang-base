package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/droomlab/drm-coupon/pkg/app/handlers"
	"github.com/droomlab/drm-coupon/pkg/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmsql"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// Run the http server
func run() error {

	conf, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "Config Initialize")
	}

	log, err := drmlog.NewZeroLogger(drmlog.Config{
		LogConfig: conf.Log,
	})
	if err != nil {
		return errors.Wrap(err, "Log Initialize")
	}

	db, err := drmsql.GetDB(drmsql.Config{
		SQLConfig: conf.Mysql,
		Log:       log,
	})
	if err != nil {
		return errors.Wrap(err, "SQL DB Initialize")
	}
	defer db.Close()

	// channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	app := handlers.NewHandlers(handlers.Config{
		Shutdown:  shutdown,
		Log:       log,
		AppConfig: conf,
	})

	server := http.Server{
		Addr:           ":" + fmt.Sprintf("%v", conf.HTTP.Port),
		Handler:        app,
		ReadTimeout:    time.Duration(conf.HTTP.ReadTimeout),
		WriteTimeout:   time.Duration(conf.HTTP.WriteTimeout),
		IdleTimeout:    time.Duration(conf.HTTP.IdleTimeout),
		MaxHeaderBytes: conf.HTTP.MaxHeaderMegabytes << 20,
		ErrorLog:       drmlog.NewServerLogger(log),
	}

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	go func() {
		log.Infof(context.TODO(), "app : API listening on port %v", conf.HTTP.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting server: %s", err)

	case sig := <-shutdown:
		log.Infof(context.TODO(), "app : Start shutdown | signal : %v", sig)

		// give outstanding requests a deadline for completion.
		timeout := time.Duration(conf.HTTP.ShutdownTimeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// asking listener to shutdown
		err := server.Shutdown(ctx)
		if err != nil {
			log.Infof(ctx, "app : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = server.Close()
		}

		if err != nil {
			return fmt.Errorf("app : could not stop server gracefully : %v", err)
		}
	}

	return nil
}
