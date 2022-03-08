package main

import (
	"fmt"
	"os"

	"github.com/droomlab/drm-coupon/pkg/app"
)

func main() {
	var configDir string = "config"
	if err := app.Run(configDir); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
