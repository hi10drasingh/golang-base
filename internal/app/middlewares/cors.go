package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/droomlab/drm-coupon/internal/app"
)

// CORS will allow access origin with "*.droom.in".
func CORS() app.Middleware {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"

	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if origin := r.Header.Get("Origin"); strings.Contains(origin, ".droom.in") {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
				w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			}

			return next(ctx, w, r)
		}

		return handler
	}

	return middleware
}
