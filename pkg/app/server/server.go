package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/app/routes"
	"github.com/droomlab/drm-coupon/pkg/appcontext"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type Server struct {
	appCtx *appcontext.AppContext
	router *httprouter.Router
}

func NewServer() (*Server, error) {
	appContext, err := appcontext.InitilizeAppContext()
	if err != nil {
		return nil, errors.Wrap(err, "AppCtx Setup Error")
	}

	router := routes.Register()

	s := &Server{appContext, router}

	return s, nil
}

func (s *Server) Serve() error {
	port := fmt.Sprintf("%v", s.appCtx.Config.Port)

	l, err := net.Listen("tcp", "localhost:"+port)

	if err != nil {
		return err
	}

	s.appCtx.Log.Info(context.TODO(), "Started listing at port "+port)

	return http.Serve(l, s.router)
}
