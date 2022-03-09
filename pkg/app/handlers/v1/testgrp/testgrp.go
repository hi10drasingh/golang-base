package testgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/drmlog"
)

// Handlers holds local dependencies for handler
type Handlers struct {
	Slug string
	Log  drmlog.Logger
}

// NewHandlers return new instance of handler
func NewHandlers(log drmlog.Logger) *Handlers {
	return &Handlers{
		Slug: "test",
		Log:  log,
	}
}

// Hello is sample route handler
func (h *Handlers) Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "hello, world!\n")
	h.Log.Info(ctx, "Workign")
	return nil
}
