package main

import (
	"fmt"
	"net/http"

	"github.com/droomlab/drm-coupon/pkg/appcontext"
	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {

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
