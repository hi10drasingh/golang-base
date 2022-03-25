package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const maxBodySize = 1048576 // 1 MB

var (
	errContentTypeNotSupported = errors.New("Content-Type header is not application/json")
	errMalformedJSON           = errors.New("Request body contains badly-formed JSON")
	errInvalidValue            = errors.New("Request body contains an invalid value")
	errUnkownField             = errors.New("Request body contains unknown field")
	errEmptyBody               = errors.New("Request body must not be empty")
	errLargeBodySize           = errors.New("Request body must not be larger than 1MB")
	errMultipleJSONObject      = errors.New("Request body must only contain a single JSON object")
)

type malformedRequestError struct {
	status int
	err    error
}

func (mr *malformedRequestError) Error() string {
	return mr.err.Error()
}

// Decode reads the body of an HTTP request looking for a JSON document.
// The body is decoded into the provided value.
func Decode(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value := r.Header.Get("Content-Type")
		if value != "application/json" {
			return &malformedRequestError{
				status: http.StatusUnsupportedMediaType,
				err:    errContentTypeNotSupported,
			}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		return decodeErr(err)
	}

	err = dec.Decode(&struct{}{})
	if errors.Is(err, io.EOF) {
		return &malformedRequestError{
			status: http.StatusBadRequest,
			err:    errMultipleJSONObject,
		}
	}

	return nil
}

func decodeErr(err error) error {
	var syntaxError *json.SyntaxError

	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		custerr := fmt.Errorf("%w (at position %d)",
			errMalformedJSON,
			syntaxError.Offset,
		)

		return &malformedRequestError{
			status: http.StatusBadRequest,
			err:    custerr,
		}

	case errors.Is(err, io.ErrUnexpectedEOF):
		return &malformedRequestError{
			status: http.StatusBadRequest,
			err:    errMalformedJSON,
		}

	case errors.As(err, &unmarshalTypeError):
		custerr := fmt.Errorf("%w %q field (at position %d)",
			errInvalidValue,
			unmarshalTypeError.Field,
			unmarshalTypeError.Offset,
		)

		return &malformedRequestError{
			status: http.StatusBadRequest,
			err:    custerr,
		}

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		custerr := fmt.Errorf("%w %s", errUnkownField, fieldName)

		return &malformedRequestError{
			status: http.StatusBadRequest,
			err:    custerr,
		}

	case errors.Is(err, io.EOF):
		return &malformedRequestError{
			status: http.StatusBadRequest,
			err:    errEmptyBody,
		}

	case err.Error() == "http: request body too large":
		return &malformedRequestError{
			status: http.StatusRequestEntityTooLarge,
			err:    errLargeBodySize,
		}

	default:
		return err
	}
}
