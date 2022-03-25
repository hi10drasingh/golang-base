package drmcontext

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/xid"
)

type ctxKey string

const ctxkey ctxKey = "golbal"

var errValueMissing = errors.New("value missing from context")

// Values holds global context appended values.
type Values struct {
	Now        time.Time
	StatusCode int
	RequestID  string
}

func SetValues(ctx context.Context) context.Context {
	requestID := xid.New().String()

	values := Values{
		Now:        time.Now().UTC(),
		StatusCode: http.StatusAccepted,
		RequestID:  requestID,
	}

	return context.WithValue(ctx, ctxkey, &values)
}

// GetValues return global value from context.
func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(ctxkey).(*Values)
	if !ok {
		return nil, errValueMissing
	}

	return v, nil
}

// SetStatusCode set value of status code in ctx.
func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(ctxkey).(*Values)
	if !ok {
		return errValueMissing
	}

	v.StatusCode = statusCode

	return nil
}
