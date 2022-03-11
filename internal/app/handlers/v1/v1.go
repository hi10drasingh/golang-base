package v1

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/droomlab/drm-coupon/internal/app"
	"github.com/droomlab/drm-coupon/internal/app/handlers/v1/testgrp"
	"github.com/droomlab/drm-coupon/internal/app/middlewares"
	drmerror "github.com/droomlab/drm-coupon/internal/app/response/error"
	drmsucces "github.com/droomlab/drm-coupon/internal/app/response/success"
	"github.com/droomlab/drm-coupon/internal/app/server"
)

// Routes register routes.
func Routes(globalHandl *app.App, deps *server.Dependencies) {
	const version = "v1"

	tgh := testgrp.NewHandlers(deps)
	group := version + "/" + tgh.Slug
	globalHandl.Handle(group, "/hello", tgh.Hello, middlewares.CheckMethod(http.MethodGet))
	globalHandl.Handle(group, "/post", tgh.Hello, middlewares.CheckMethod(http.MethodPost))

	// NotFound Handler
	globalHandl.Handle("", "/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		if r.URL.Path == "/" {
			data := [][]byte{}

			err = app.Respond(ctx, w, drmsucces.NewResponse(data, "Wellcome!"), http.StatusOK)

			return errors.Wrap(err, "Handler Respond")
		}

		status := drmerror.StatusNotFound
		message := drmerror.StatusNotFoundMsg

		err = errors.New(message)

		return drmerror.NewRequestError(err, status, message)
	})
}
