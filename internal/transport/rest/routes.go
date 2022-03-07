package rest

import "net/http"

// globalMiddleware setup
func (h *Handlers) globalMiddlewares(method string, handler drmHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var final http.HandlerFunc

		switch method {
		case http.MethodGet:
			final = h.errorHandler(h.get(h.panicRecovery(h.logger(h.cors(handler)))))
		case http.MethodPost:
			final = h.errorHandler(h.post(h.panicRecovery(h.logger(h.cors(handler)))))
		}

		final(w, r)
	}
}

func (h *Handlers) setupRoutes(router *http.ServeMux) {
	router.Handle("/hello", h.globalMiddlewares(http.MethodGet, h.hello()))
	router.Handle("/post", h.globalMiddlewares(http.MethodPost, h.hello()))
}
