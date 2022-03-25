package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/droomlab/drm-coupon/internal/app"
	drmerror "github.com/droomlab/drm-coupon/internal/app/response/error"
)

var errMethodNotAllowed = errors.New("request method not allowed")

// CheckMethod provide reuqest method checking for
// default ServeMux handler.
func CheckMethod(method string) app.Middleware {
	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if method != r.Method {
				code := drmerror.StatusMethodNotAllowed
				msg := drmerror.StatusMethodNotAllowedMsg

				return drmerror.NewRequestError(errMethodNotAllowed, code, msg)
			}

			return next(ctx, w, r)
		}

		return handler
	}

	return middleware
}
