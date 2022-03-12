package server

import (
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/drmlog"
)

// New return new http server with HTTP config.
func New(handler http.Handler, deps *Dependencies) *http.Server {
	return &http.Server{
		Addr:           ":" + fmt.Sprintf("%v", deps.Config.HTTP.Port),
		Handler:        handler,
		ReadTimeout:    deps.Config.HTTP.ReadTimeout.Time,
		WriteTimeout:   deps.Config.HTTP.WriteTimeout.Time,
		IdleTimeout:    deps.Config.HTTP.IdleTimeout.Time,
		MaxHeaderBytes: deps.Config.HTTP.MaxHeaderMegabytes << 20,
		ErrorLog:       drmlog.NewServerLogger(deps.Log),
	}
}
