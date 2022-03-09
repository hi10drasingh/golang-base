package app

import (
	"context"
	"errors"
	"time"
)

type ctxKey int

const ctxkey ctxKey = 1

// Values holds global context appended values
type Values struct {
	Now        time.Time
	StatusCode int
}

// GetValues return global value from context
func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(ctxkey).(*Values)
	if !ok {
		return nil, errors.New("value missing from context")
	}
	return v, nil
}

// SetStatusCode set value of status code in ctx
func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(ctxkey).(*Values)
	if !ok {
		return errors.New("value missing from context")
	}
	v.StatusCode = statusCode
	return nil
}
