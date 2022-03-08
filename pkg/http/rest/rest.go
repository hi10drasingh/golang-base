package rest

import (
	"fmt"
	"net/http"
	"time"
)

// GetServer creates new instance of HTTP server
func (h *Handlers) GetServer() *http.Server {

	router := http.NewServeMux()
	h.setupRoutes(router)

	server := http.Server{
		Addr:           ":" + fmt.Sprintf("%v", h.AppCtx.Config.HTTP.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(h.AppCtx.Config.HTTP.ReadTimeout),
		WriteTimeout:   time.Duration(h.AppCtx.Config.HTTP.WriteTimeout),
		IdleTimeout:    time.Duration(h.AppCtx.Config.HTTP.IdleTimeout),
		MaxHeaderBytes: h.AppCtx.Config.HTTP.MaxHeaderMegabytes << 20,
	}

	return &server
}
