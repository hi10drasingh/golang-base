package middlewares

import (
	"context"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/app"
	"github.com/droomlab/drm-coupon/pkg/drmerrors"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
)

// Errors provide error handling for router handlers
func Errors(log drmlog.Logger) app.Middleware {
	m := func(handler app.Handler) app.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(err, "Error in Handler")

				var er drmerrors.ErrorResponse

				switch {
				case drmerrors.IsRequestError(err):
					reqErr := drmerrors.GetRequestError(err)
					er = drmerrors.ErrorResponse{
						Code:       reqErr.Code,
						Message:    reqErr.Message,
						StatusCode: reqErr.Status,
					}

				default:
					er = drmerrors.ErrorResponse{
						Code:       drmerrors.CodeFailed,
						Message:    http.StatusText(http.StatusInternalServerError),
						StatusCode: http.StatusInternalServerError,
					}
				}

				// Respond with the error back to the client.
				if err := app.Respond(ctx, w, er, er.StatusCode); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if ok := app.IsShutdown(err); ok {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
