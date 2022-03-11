package middlewares

import (
	"context"
	"net/http"

	"github.com/droomlab/drm-coupon/internal/app"
	drmerror "github.com/droomlab/drm-coupon/internal/app/response/error"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/pkg/errors"
)

// Errors provide error handling for router handlers.
func Errors(log drmlog.Logger) app.Middleware {
	middleware := func(next app.Handler) app.Handler {
		handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := next(ctx, w, r); err != nil {
				log.Error(ctx, err, "Error in Handler")

				var apperr drmerror.ErrorResponse

				switch {
				case drmerror.IsRequestError(err):
					reqErr := drmerror.GetRequestError(err)
					apperr = drmerror.ErrorResponse{
						Code:       reqErr.Code,
						Message:    reqErr.Message,
						StatusCode: reqErr.Status,
					}

				default:
					apperr = drmerror.NewErrorResponse(
						drmerror.StatusInternalServerError,
						drmerror.StatusInternalServerErrorMsg,
					)
				}

				// Respond with the error back to the client.
				if respondErr := app.Respond(ctx, w, apperr, apperr.StatusCode); respondErr != nil {
					return errors.Wrap(respondErr, "App Respond")
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if ok := app.IsShutdown(err); ok {
					return errors.Wrap(err, "App Shutdown")
				}
			}

			return nil
		}

		return handler
	}

	return middleware
}
