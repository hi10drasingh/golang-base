package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/justinas/alice"

	"github.com/droomlab/drm-coupon/internal/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// RequestLogger provided logging of all client request.
func RequestLogger(log zerolog.Logger) app.Middleware {
	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			c := alice.New()

			// Install the logger handler with default output on the console
			c = c.Append(hlog.NewHandler(log))

			// Install some provided extra handler to set some request's context fields.
			// Thanks to that handler, all our logs will come with some prepopulated fields.
			c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
				hlog.FromRequest(r).Info().
					Str("method", r.Method).
					Stringer("url", r.URL).
					Int("status", status).
					Int("size", size).
					Dur("duration", duration).
					Msg("")
			}))
			c = c.Append(hlog.RemoteAddrHandler("ip"))
			c = c.Append(hlog.UserAgentHandler("user_agent"))
			c = c.Append(hlog.RefererHandler("referer"))
			c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

			var err error

			// Here is your final handler
			c.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hlog.FromRequest(r).Info().
					Str("user", "current user").
					Str("status", "ok").
					Msg("Something happened")
				err = next(ctx, w, r)
			})).ServeHTTP(w, r)

			return err
		}

		return handler
	}

	return middleware
}
