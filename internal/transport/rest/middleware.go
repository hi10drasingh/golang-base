package rest

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

func (h *Handlers) logger(han http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			h.AppCtx.Log.Infof("The client %s requested %v \n", r.RemoteAddr, r.URL)
			h.AppCtx.Log.Infof("It took %s to serve the request \n", time.Since(startTime))
		}()
		han(w, r)
	}
}

func (h *Handlers) cors(han http.HandlerFunc) http.HandlerFunc {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); strings.Contains(origin, ".droom.in") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		}

		han(w, r)
	}
}

func (h *Handlers) panicRecovery(han http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				var err error

				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}

				if err != nil {
					h.AppCtx.Log.Error(err, "panic error occurred")
				}
			}
		}()
		han(w, r)
	}
}

// global middleware setup
func (h *Handlers) setUpGlobalMiddlewares(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		final := h.panicRecovery(h.logger(h.cors(handler.ServeHTTP)))
		final(w, r)
	})
}
