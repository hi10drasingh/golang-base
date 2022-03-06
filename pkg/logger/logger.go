package logger

type LogConfig struct {
	Dir   string `json:"dir"`
	Level int    `json:"level"`
}

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
