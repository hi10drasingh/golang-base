package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/droomlab/drm-coupon/internal/appcontext"
)

type Handlers struct {
	AppCtx *appcontext.AppContext
}

func NewHandlers(AppCtx *appcontext.AppContext) *Handlers {
	return &Handlers{
		AppCtx: AppCtx,
	}
}

func (h *Handlers) GetServer() *http.Server {

	router := http.NewServeMux()
	h.SetupRoutes(router)

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
