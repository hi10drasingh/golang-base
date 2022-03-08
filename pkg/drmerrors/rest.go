package drmerrors

import (
	"errors"
	"net/http"
)

var (
	// NotFound 404
	NotFound string = http.StatusText(http.StatusNotFound)
	// MethodNotAllowed 405 Method Not Allowed
	MethodNotAllowed string = http.StatusText(http.StatusMethodNotAllowed)
	// InternalServerError 500 Internal Server Error
	InternalServerError string = http.StatusText(http.StatusInternalServerError)
)

// CodeFailed used for request errors
const CodeFailed = "failed"

// ErrorResponse prodives sending error response to client
type ErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// RequestError contains details of Req Error
type RequestError struct {
	Err     error
	Code    string
	Message string
	Status  int
}

// NewRequestError return a new HTTP error with provided err, status and code
func NewRequestError(err error, status int, message string) error {
	return &RequestError{
		Err:     err,
		Code:    "failed",
		Message: message,
		Status:  status,
	}
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}

// IsRequestError check if a err is of type RequestError
func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

// GetRequestError Converts err into Type RequestError
func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
