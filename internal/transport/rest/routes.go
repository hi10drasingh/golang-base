package rest

import "net/http"

func (h *Handlers) SetupRoutes(router *http.ServeMux) {
	router.HandleFunc("/hello", h.logger(h.hello()))
}
