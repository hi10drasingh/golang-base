package handlers

import (
	"net/http"
	"os"

	"github.com/droomlab/drm-coupon/internal/app"
	"github.com/droomlab/drm-coupon/internal/app/dependency"
	v1 "github.com/droomlab/drm-coupon/internal/app/handlers/v1"
	"github.com/droomlab/drm-coupon/internal/app/middlewares"
)

// Config holds global dependencies for handler.
type Config struct {
	Shutdown chan os.Signal
	Deps     *dependency.Dependency
}

// NewHandlers return new instance of handler.
func NewHandlers(conf Config) http.Handler {
	globalHandl := app.NewApp(
		conf.Shutdown,
		conf.Deps.Log,
		conf.Deps.Config,
		middlewares.RequestLogger(conf.Deps.RequestLog),
		middlewares.CORS(),
		middlewares.Errors(conf.Deps.Log),
		middlewares.Recovery(conf.Deps.Log),
	)

	v1.Routes(globalHandl, conf.Deps)

	return globalHandl
}
