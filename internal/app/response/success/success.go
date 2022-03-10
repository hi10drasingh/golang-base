package success

import "net/http"

const code = "success"

// Response hold succes response data
type Response struct {
	Code       string      `json:"code"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
}

// NewResponse returns new Response struct
// which will contain data and message string
func NewResponse(data interface{}, message string) Response {
	return Response{
		Code:       code,
		Data:       data,
		Message:    message,
		StatusCode: http.StatusOK,
	}
}
