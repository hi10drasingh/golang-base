package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/drmcontext"
	"github.com/pkg/errors"
)

const requestIDHeaderName = "Request-ID"

// Respond send response to client with provide statusCode and Data.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	// Set the status code for the request logger middleware.
	err := drmcontext.SetStatusCode(ctx, statusCode)
	if err != nil {
		return errors.Wrap(err, "Respond Setting Context Status Code")
	}

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)

		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "Respond Data Marshal")
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	_, err = w.Write(jsonData)

	if err != nil {
		return errors.Wrap(err, "Respond Write JSON")
	}

	return nil
}

func SetRequestIDHeader(ctx context.Context, w http.ResponseWriter) {
	drmctx, _ := drmcontext.GetValues(ctx)

	// Setting Request Id In Headers
	w.Header().Set(requestIDHeaderName, drmctx.RequestID)
}
