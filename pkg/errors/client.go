package errors

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
