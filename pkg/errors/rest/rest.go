package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	// MethodNotAllowed 405 Method Not Allowed
	MethodNotAllowed string = http.StatusText(http.StatusMethodNotAllowed)
	// InternalServerError 500 Internal Server Error
	InternalServerError string = http.StatusText(http.StatusInternalServerError)
)

// Error Contains Error , Message , Code And Status
type Error struct {
	Cause   error  `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"statusCode"`
}

// Error returns message attached with error
func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

// ErrorObj return actual error obj
func (e *Error) ErrorObj() error {
	return e.Cause
}

// ResponseBody returns JSON response body.
func (e *Error) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

// ResponseHeaders returns http status code and headers.
func (e *Error) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

// NewError return a new HTTP error with provided err, status and code
func NewError(err error, status int, message string) error {
	return &Error{
		Cause:   err,
		Code:    "failed",
		Message: message,
		Status:  status,
	}
}
