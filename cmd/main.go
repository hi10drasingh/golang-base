package main

import (
	"fmt"
	"os"

	"github.com/droomlab/drm-coupon/pkg/app/server"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s, err := server.NewServer()

	if err != nil {
		return errors.Wrap(err, "Server Create")
	}

	return s.Serve()

}
