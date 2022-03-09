package testgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"
	"github.com/tsenart/nap"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handlers holds local dependencies for handler
type Handlers struct {
	Slug  string
	Log   drmlog.Logger
	SQL   *nap.DB
	NoSQL *mongo.Client
	RMQ   *drmrmq.RabbitMQ
}

// NewHandlers return new instance of handler
func NewHandlers(handlers *Handlers) *Handlers {
	handlers.Slug = "test"
	return handlers
}

// Hello is sample route handler
func (h Handlers) Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "hello, world!\n")
	h.Log.Info(ctx, "Workign")
	return nil
}
