package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/internal/appcontext"
)

type Server struct {
	appCtx *appcontext.AppContext
	httpServer *http.Server
}

func NewServer(appCtx *appcontext.AppContext, handler http.Handler) (*Server) {
	return &Server{
		appCtx: appCtx,
		httpServer: &http.Server{
			Addr:           ":" + fmt.Sprintf("%v", appCtx.Config.HTTP.Port),
			Handler:        handler,
			ReadTimeout:    appCtx.Config.HTTP.ReadTimeout,
			WriteTimeout:   appCtx.Config.HTTP.WriteTimeout,
			MaxHeaderBytes: appCtx.Config.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
