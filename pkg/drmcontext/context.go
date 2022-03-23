package drmcontext

import (
	"context"
	"errors"

	"github.com/rs/zerolog/hlog"
)

// ReqIDKey will be the key name of request id in log msg.
const ReqIDKey = "req_id"

var errNoReqID = errors.New("context value of request id missing")

func RequestIDCtx(ctx context.Context) (string, error) {
	value, ok := hlog.IDFromCtx(ctx)

	if !ok {
		return "", errNoReqID
	}

	return value.String(), nil
}
