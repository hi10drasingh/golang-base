package server

import (
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/appcontext"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	appCtx *appcontext.AppContext
}

func Serve() {
	appContext, err := appcontext.InitilizeAppContext()
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/hello/:name", hello)

	appContext.Log.Info("Server started at port 8888")
	appContext.Log.Fatal(http.ListenAndServe(":8888", router), "Server Failure")
}