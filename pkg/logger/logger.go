package logger

import "context"

var Log Logger

type Logger interface {
	Errorf(ctx context.Context, err error, format string, args ...interface{})
	Error(ctx context.Context, err error, msg string)
	Fatalf(ctx context.Context, err error, format string, args ...interface{})
	Fatal(ctx context.Context, err error, msg string)
	Infof(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, msg string)
	Debugf(ctx context.Context, format string, args ...interface{})
	Debug(ctx context.Context, msg string)
}
