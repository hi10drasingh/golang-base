package app

import (
	"context"
	"errors"
	"time"
)

type ctxKey int

const ctxkey ctxKey = 1

var errValueMissing = errors.New("value missing from context")

// Values holds global context appended values.
type Values struct {
	Now        time.Time
	StatusCode int
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
