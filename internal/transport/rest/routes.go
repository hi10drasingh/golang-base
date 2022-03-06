package rest

import "net/http"

func (h *Handlers) setupRoutes(router *http.ServeMux) {
	router.HandleFunc("/hello", get(h.hello()))
	router.HandleFunc("/post", post(h.hello()))
}

func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, http.MethodGet)
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, http.MethodPost)
}
