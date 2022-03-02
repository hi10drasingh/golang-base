package logger

var Log Logger

type Logger interface {
	Errorf(err error, msg string, args ...interface{})
	Error(err error, msg string)
	Fatalf(err error, msg string, args ...interface{})
	Fatal(err error)
	Infof(msg string, args ...interface{})
	Info(msg string)
	Debugf(msg string, args ...interface{})
	Debug(msg string)
}

