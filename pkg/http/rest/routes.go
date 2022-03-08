package rest

import (
	"net/http"
)

type route struct {
	pattern string
	method  string
	handler drmHandler
}

// globalMiddleware setup
func (h *Handlers) globalMiddlewares(router *http.ServeMux, routes []route) {
	for i := 0; i < len(routes); i++ {
		newRoute := routes[i]

		router.HandleFunc(newRoute.pattern, func(w http.ResponseWriter, r *http.Request) {
			var final http.HandlerFunc

			switch newRoute.method {
			case http.MethodGet:
				final = h.errorHandler(h.get(h.panicRecovery(h.logger(h.cors(newRoute.handler)))))
			case http.MethodPost:
				final = h.errorHandler(h.post(h.panicRecovery(h.logger(h.cors(newRoute.handler)))))
			}

			final(w, r)
		})
	}
}

func (h *Handlers) setupRoutes(router *http.ServeMux) {
	routes := []route{
		{"/hello", http.MethodGet, h.hello()},
		{"/post", http.MethodPost, h.hello()},
	}

	h.globalMiddlewares(router, routes)
}
