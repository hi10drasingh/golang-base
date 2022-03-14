package testgrp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/internal/app/dependency"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"
	"github.com/tsenart/nap"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handlers holds local dependencies for handler.
type Handlers struct {
	Slug  string
	Log   drmlog.Logger
	SQL   *nap.DB
	NoSQL *mongo.Client
	RMQ   *drmrmq.RabbitMQ
}

// NewHandlers return new instance of handler.
func NewHandlers(deps *dependency.Dependency) *Handlers {
	return &Handlers{
		Slug:  "test",
		Log:   deps.Log,
		SQL:   deps.SQL,
		NoSQL: deps.NoSQL,
		RMQ:   deps.RMQ,
	}
}

// Hello is sample route handler.
func (h Handlers) Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "hello, world!\n")

	return nil
}
