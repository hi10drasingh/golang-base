package errors

import (
	"errors"
	"net/http"
)

const code = "failed"

const (

	// StatusNotFound 404 Route Not Found
	StatusNotFound int = http.StatusNotFound

	// StatusMethodNotAllowed 405 Method Not Allowed
	StatusMethodNotAllowed int = http.StatusMethodNotAllowed

	// StatusInternalServerError 500 Internal Server Error
	StatusInternalServerError int = http.StatusInternalServerError
)

var statusText = map[int]string{
	StatusNotFound:            "Not Found",
	StatusMethodNotAllowed:    "Method Not Allowed",
	StatusInternalServerError: "Internal Server Error",
}

// StatusText return message crossponding to a status code
func StatusText(code int) string {
	return statusText[code]
}

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
		Code:    code,
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

// NewErrorResponse return a new error reponse with provided status and code
func NewErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{
		Code:       code,
		Message:    message,
		StatusCode: status,
	}
}
