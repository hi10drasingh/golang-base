package logger

import "context"

type LogConfig struct {
	Dir   string `json:"dir"`
	Level int    `json:"level"`
}

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
