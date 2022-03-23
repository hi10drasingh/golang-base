package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/justinas/alice"

	"github.com/droomlab/drm-coupon/internal/app"
	"github.com/droomlab/drm-coupon/pkg/drmcontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

const (
	method           = "method"
	url              = "url"
	status           = "status"
	size             = "size"
	duration         = "duration"
	ipAddr           = "ip"
	userAgent        = "user_agent"
	referer          = "referer"
	requestHeaderKey = "Request-Id"
)

// RequestLogger provided logging of all client request.
func RequestLogger(log *zerolog.Logger) app.Middleware {
	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			club := alice.New()

			// Install the logger handler with default output on the console
			club = club.Append(hlog.NewHandler(*log))

			// Install some provided extra handler to set some request's context fields.
			// Thanks to that handler, all our logs will come with some prepopulated fields.
			club = club.Append(hlog.AccessHandler(func(r *http.Request, statusCode, bodySize int, reqTime time.Duration) {
				hlog.FromRequest(r).Info().
					Str(method, r.Method).
					Stringer(url, r.URL).
					Int(status, statusCode).
					Int(size, bodySize).
					Dur(duration, reqTime).
					Msg("")
			}))
			club = club.Append(hlog.RemoteAddrHandler(ipAddr))
			club = club.Append(hlog.UserAgentHandler(userAgent))
			club = club.Append(hlog.RefererHandler(referer))
			club = club.Append(hlog.RequestIDHandler(drmcontext.ReqIDKey, requestHeaderKey))

			var err error

			// Here is your final handler
			club.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx = r.Context()
				err = next(ctx, w, r)
			})).ServeHTTP(w, r)

			return err
		}

		return handler
	}

	return middleware
}
