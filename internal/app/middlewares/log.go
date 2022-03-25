package middlewares

import (
	"context"
	"io"
	"net/http"

	"github.com/droomlab/drm-coupon/internal/app"
	"github.com/droomlab/drm-coupon/pkg/drmcontext"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	method    = "method"
	url       = "url"
	body      = "body"
	requestID = "request_id"
)

// ContextLogger provided logging of all client request.
func ContextLogger(log drmlog.Logger) app.Middleware {
	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			defaultLogger := log.GetLogger()

			newLogger := defaultLogger.With().Logger()

			ctx = newLogger.WithContext(ctx)

			r = r.WithContext(ctx)

			drmctx, _ := drmcontext.GetValues(ctx)

			contextLogger := zerolog.Ctx(ctx)

			reqBody, err := io.ReadAll(r.Body)
			defer r.Body.Close()

			if err != nil {
				return errors.Wrap(err, "Request Body Read")
			}

			contextLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
				// Adding method in log
				context := c.Str(requestID, drmctx.RequestID)

				// Adding method in log
				context = context.Str(method, r.Method)

				// Adding url in log
				context = context.Stringer(url, r.URL)

				// Adding rquest body
				context = context.RawJSON(body, reqBody)

				return context
			})

			return next(ctx, w, r)
		}

		return handler
	}

	return middleware
}
