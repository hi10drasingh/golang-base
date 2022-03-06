package rest

import (
	"fmt"
	"net/http"
)

func (h *Handlers) hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello, world!\n")
	}
}
