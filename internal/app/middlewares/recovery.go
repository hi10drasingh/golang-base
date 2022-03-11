package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/droomlab/drm-coupon/pkg/drmlog"

	"github.com/droomlab/drm-coupon/internal/app"
)

// Recovery provides panic handling for application.
func Recovery(log drmlog.Logger) app.Middleware {
	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()

					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
				}
			}()

			err = next(ctx, w, r)

			return err
		}

		return handler
	}

	return middleware
}
