package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/app"

	"github.com/droomlab/drm-coupon/pkg/drmerrors"
)

// CheckMethod provide reuqest method checking for
// default ServeMux handler
func CheckMethod(method string) app.Middleware {
	m := func(handler app.Handler) app.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if method != r.Method {
				return drmerrors.NewRequestError(errors.New("testing error method"), http.StatusMethodNotAllowed, drmerrors.MethodNotAllowed)
			}
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
