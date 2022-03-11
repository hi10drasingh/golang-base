package app

import "errors"

type shutdownError struct {
	Message string
}

// NewShutdownError generates new shutdown error with provided msg.
func NewShutdownError(message string) error {
	return &shutdownError{message}
}

func (se *shutdownError) Error() string {
	return se.Message
}

// IsShutdown checks if err is of type shutdownError.
func IsShutdown(err error) bool {
	var se *shutdownError

	return errors.As(err, &se)
}
