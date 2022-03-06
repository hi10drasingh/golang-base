package rest

import "net/http"

func (h *Handlers) setupRoutes(router *http.ServeMux) {
	router.HandleFunc("/hello", h.hello())
}
