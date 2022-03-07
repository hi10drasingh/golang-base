package drmerrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ClientError is an error whose details to be shared with client.
type ClientError interface {
	// Error returns message attached with error
	Error() string
	// ErrorObj return actual error obj
	ErrorObj() error
	// ResponseBody returns response body.
	ResponseBody() ([]byte, error)
	// ResponseHeaders returns http status code and headers.
	ResponseHeaders() (int, map[string]string)
}

var (
	// MethodNotAllowed 405 Method Not Allowed
	MethodNotAllowed string = http.StatusText(http.StatusMethodNotAllowed)
	// InternalServerError 500 Internal Server Error
	InternalServerError string = http.StatusText(http.StatusInternalServerError)
)

// HTTPError Contains Error , Message , Code And Status
type HTTPError struct {
	Cause   error  `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"statusCode"`
}

// Error returns message attached with error
func (e *HTTPError) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

// ErrorObj return actual error obj
func (e *HTTPError) ErrorObj() error {
	return e.Cause
}

// ResponseBody returns JSON response body.
func (e *HTTPError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

// ResponseHeaders returns http status code and headers.
func (e *HTTPError) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

// NewHTTPError return a new HTTP error with provided err, status and code
func NewHTTPError(err error, status int, message string) error {
	return &HTTPError{
		Cause:   err,
		Code:    "failed",
		Message: message,
		Status:  status,
	}
}
