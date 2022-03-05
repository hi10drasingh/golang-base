package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Register() *httprouter.Router {
	router := httprouter.New()

	router.GET("/hello/:name", hello)

	return router
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
