package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/droomlab/drm-coupon/internal/appcontext"
)

// Handlers Struct is responsible for server creation and registration of route-handlers mapping
type Handlers struct {
	AppCtx *appcontext.AppContext
}

// NewHandlers returns new instance of Handles Struct
func NewHandlers(appCtx *appcontext.AppContext) *Handlers {
	return &Handlers{
		AppCtx: appCtx,
	}
}

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
