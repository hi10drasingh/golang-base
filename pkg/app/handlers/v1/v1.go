package v1

import (
	"net/http"

	"github.com/droomlab/drm-coupon/domain/middlewares"
	"github.com/droomlab/drm-coupon/pkg/app"
	"github.com/droomlab/drm-coupon/pkg/app/handlers/v1/testgrp"
	"github.com/droomlab/drm-coupon/pkg/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
)

// Config holds dependencies for Router
type Config struct {
	Log       drmlog.Logger
	AppConfig *config.App
}

// Routes register routes
func Routes(a *app.App, conf Config) {
	const version = "v1"

	tgh := testgrp.NewHandlers(
		conf.Log,
	)
	group := version + "/" + tgh.Slug
	a.Handle(group, "/hello", tgh.Hello, middlewares.CheckMethod(http.MethodGet))
	a.Handle(group, "/post", tgh.Hello, middlewares.CheckMethod(http.MethodPost))
}
