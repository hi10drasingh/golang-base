package rest

import (
	"errors"
	"net/http"
	"strings"
	"time"

	drmerrors "github.com/droomlab/drm-coupon/pkg/errors"
)

func (h *Handlers) logger(han drmHandler) drmHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		startTime := time.Now()
		defer func() {
			h.AppCtx.Log.Infof("The client %s requested %v \n", r.RemoteAddr, r.URL)
			h.AppCtx.Log.Infof("It took %s to serve the request \n", time.Since(startTime))
		}()
		return han(w, r)
	}
}

func (h *Handlers) cors(han drmHandler) drmHandler {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	return func(w http.ResponseWriter, r *http.Request) error {
		if origin := r.Header.Get("Origin"); strings.Contains(origin, ".droom.in") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		}

		return han(w, r)
	}
}

func (h *Handlers) panicRecovery(han drmHandler) drmHandler {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}
			}
		}()
		err = han(w, r)
		return err
	}
}

func (h *Handlers) checkMethod(dh drmHandler, method string) drmHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if method != r.Method {
			return drmerrors.NewHTTPError(errors.New("failedcnkcn"), http.StatusMethodNotAllowed, drmerrors.MethodNotAllowed)
		}
		return dh(w, r)
	}
}

func (h *Handlers) get(dh drmHandler) drmHandler {
	return h.checkMethod(dh, http.MethodGet)
}

func (h *Handlers) post(dh drmHandler) drmHandler {
	return h.checkMethod(dh, http.MethodPost)
}

func (h *Handlers) errorHandler(han drmHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error = han(w, r) // Call handler function
		if err == nil {
			return
		}

		// Check if it is a ClientError.
		clientError, ok := err.(drmerrors.ClientError)

		if !ok {
			h.AppCtx.Log.Error(err, "Server Error")
			// If the error is not ClientError, assume that it is ServerError.
			w.WriteHeader(http.StatusInternalServerError) // return 500 Internal Server Error.
			return
		}

		body, err := clientError.ResponseBody() // Try to get response body of ClientError.
		if err != nil {
			h.AppCtx.Log.Error(err, "ClientError Decode")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		status, headers := clientError.ResponseHeaders() // Get http status code and headers.
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		_, err = w.Write(body)

		if err != nil {
			h.AppCtx.Log.Error(err, "Response Body Write")
			return
		}

		h.AppCtx.Log.Error(clientError.ErrorObj(), clientError.Error())
	}
}
