package logger

// LogConfig holds configuration for logger
type LogConfig struct {
	Dir   string `json:"dir"`
	Level int    `json:"level"`
}

// Logger is a generic interface for standard logging libraries
type Logger interface {
	Errorf(err error, format string, args ...interface{})
	Error(err error, msg string)
	Fatalf(err error, format string, args ...interface{})
	Fatal(err error, msg string)
	Infof(format string, args ...interface{})
	Info(msg string)
	Debugf(format string, args ...interface{})
	Debug(msg string)
}
