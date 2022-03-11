package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/droomlab/drm-coupon/pkg/drmlog"
)

// New return new http server with HTTP config.
func New(handler http.Handler, deps *Dependencies) *http.Server {
	return &http.Server{
		Addr:           ":" + fmt.Sprintf("%v", deps.Config.HTTP.Port),
		Handler:        handler,
		ReadTimeout:    time.Duration(deps.Config.HTTP.ReadTimeout),
		WriteTimeout:   time.Duration(deps.Config.HTTP.WriteTimeout),
		IdleTimeout:    time.Duration(deps.Config.HTTP.IdleTimeout),
		MaxHeaderBytes: deps.Config.HTTP.MaxHeaderMegabytes << 20,
		ErrorLog:       drmlog.NewServerLogger(deps.Log),
	}
}
