package rest

import (
	"net/http"
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
