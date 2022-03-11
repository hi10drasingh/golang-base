package app

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Decode reads the body of an HTTP request looking for a JSON document.
// The body is decoded into the provided value.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(val)

	return errors.Wrap(err, "Decode Request Body")
}
