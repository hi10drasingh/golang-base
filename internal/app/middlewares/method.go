package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/droomlab/drm-coupon/internal/app"

	drmerrors "github.com/droomlab/drm-coupon/internal/app/response/error"
)

// CheckMethod provide reuqest method checking for
// default ServeMux handler
func CheckMethod(method string) app.Middleware {
	m := func(handler app.Handler) app.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if method != r.Method {
				return drmerrors.NewRequestError(errors.New("testing error method"), drmerrors.StatusMethodNotAllowed, drmerrors.StatusText(drmerrors.StatusMethodNotAllowed))
			}
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
