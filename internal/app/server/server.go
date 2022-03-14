package server

import (
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmtime"
)

// Config holds configuration for  server.
type Config struct {
	// Host               string     `json:"host" validate:"required"`
	Port               int                `json:"port" validate:"required,number"`
	ReadTimeout        drmtime.CustomTime `json:"readTimeout" validate:"required"`
	WriteTimeout       drmtime.CustomTime `json:"writeTimeout" validate:"required"`
	IdleTimeout        drmtime.CustomTime `json:"idleTimeout" validate:"required"`
	ShutdownTimeout    drmtime.CustomTime `json:"shutdownTimeout" validate:"required"`
	MaxHeaderMegabytes int                `json:"maxHeaderMegaBytes" validate:"required,number"`
}

// New return new http server with HTTP config.
func New(handler http.Handler, conf *Config, log drmlog.Logger) *http.Server {
	return &http.Server{
		Addr:           ":" + fmt.Sprintf("%v", conf.Port),
		Handler:        handler,
		ReadTimeout:    conf.ReadTimeout.Time,
		WriteTimeout:   conf.WriteTimeout.Time,
		IdleTimeout:    conf.IdleTimeout.Time,
		MaxHeaderBytes: conf.MaxHeaderMegabytes << 20,
		ErrorLog:       drmlog.NewServerLogger(log),
	}
}
