package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/droomlab/drm-coupon/internal/app"
)

// Recovery provides panic handling for application
func Recovery() app.Middleware {
	m := func(handler app.Handler) app.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {

					// Stack trace will be provided.
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
				}
			}()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
