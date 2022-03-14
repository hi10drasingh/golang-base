package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/droomlab/drm-coupon/pkg/drmlog"

	"github.com/droomlab/drm-coupon/internal/app"
)

var ErrPanic = errors.New("PANIC")

// Recovery provides panic handling for application.
func Recovery(log drmlog.Logger) app.Middleware {
	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()

					err = fmt.Errorf("[%v] TRACE[%s]: %w", rec, string(trace), ErrPanic)
				}
			}()

			err = next(ctx, w, r)

			return err
		}

		return handler
	}

	return middleware
}
