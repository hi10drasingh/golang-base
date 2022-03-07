package httperr

import (
	"encoding/json"
	"fmt"
)

// HTTPError implements ClientError interface.
type HTTPError struct {
	Cause   error  `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"statusCode"`
}

func (e *HTTPError) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
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
