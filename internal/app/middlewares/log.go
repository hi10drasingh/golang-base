package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/droomlab/drm-coupon/internal/app"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
)

// Logger provided logging of all client request.
func Logger(log drmlog.Logger) app.Middleware {
	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			startTime := time.Now()
			defer func() {
				log.Infof(ctx, "The client %s requested %v \n", r.RemoteAddr, r.URL)
				log.Infof(ctx, "It took %s to serve the request \n", time.Since(startTime))
			}()

			return next(ctx, w, r)
		}

		return handler
	}

	return middleware
}
