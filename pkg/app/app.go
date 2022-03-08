package app

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/droomlab/drm-coupon/pkg/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
)

// Handler type defines type of request handlers
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App struct hold application level dependencies
type App struct {
	router   *http.ServeMux
	shutdown chan os.Signal
	log      drmlog.Logger
	config   *config.App
	mw       []Middleware
}

// NewApp returns new instance of App
func NewApp(shutdown chan os.Signal, log drmlog.Logger, config *config.App, mw ...Middleware) *App {
	router := http.NewServeMux()

	return &App{
		router:   router,
		shutdown: shutdown,
		log:      log,
		config:   config,
		mw:       mw,
	}
}

// SignalShutdown used to gracefully shutdown app
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// ServeHTTP used to make app a valid http.Handler
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// Handle used to register new routes
func (a *App) Handle(group string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)

	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		v := Values{
			Now: time.Now().UTC(),
		}

		ctx = context.WithValue(ctx, ctxkey, &v)

		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.router.HandleFunc(finalPath, h)
}
