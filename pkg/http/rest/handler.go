package rest

import (
	"fmt"
	"net/http"
)

type drmHandler func(w http.ResponseWriter, r *http.Request) error

func (h *Handlers) hello() drmHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "hello, world!\n")

		return nil
	}
}
