package rest

import (
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/appcontext"
)

// Handlers Struct is responsible for server creation and registration of route-handlers mapping
type Handlers struct {
	AppCtx *appcontext.AppContext
}

// drmHandler will be used to return error from handlers
type drmHandler func(w http.ResponseWriter, r *http.Request) error

// NewHandlers returns new instance of Handles Struct
func NewHandlers(appCtx *appcontext.AppContext) *Handlers {
	return &Handlers{
		AppCtx: appCtx,
	}
}

func (h *Handlers) hello() drmHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "hello, world!\n")

		return nil
	}
}
