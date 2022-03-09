package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/droomlab/drm-coupon/domain/middlewares"
	"github.com/droomlab/drm-coupon/internal/app"
	"github.com/droomlab/drm-coupon/internal/app/handlers/v1/testgrp"
	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/drmerrors"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"
	"github.com/tsenart/nap"
	"go.mongodb.org/mongo-driver/mongo"
)

// Config holds dependencies for Router
type Config struct {
	Log       drmlog.Logger
	AppConfig *config.App
	SQL       *nap.DB
	NoSQL     *mongo.Client
	RMQ       *drmrmq.RabbitMQ
}

// Routes register routes
func Routes(a *app.App, conf Config) {
	const version = "v1"

	tgh := testgrp.NewHandlers(&testgrp.Handlers{
		Log:   conf.Log,
		SQL:   conf.SQL,
		NoSQL: conf.NoSQL,
		RMQ:   conf.RMQ,
	})
	group := version + "/" + tgh.Slug
	a.Handle(group, "/hello", tgh.Hello, middlewares.CheckMethod(http.MethodGet))
	a.Handle(group, "/post", tgh.Hello, middlewares.CheckMethod(http.MethodPost))

	// NotFound Handler
	a.Handle("", "/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return drmerrors.NewRequestError(errors.New(drmerrors.NotFound), http.StatusNotFound, drmerrors.NotFound)
	})
}
