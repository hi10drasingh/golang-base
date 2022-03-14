package drmlog

import (
	"context"
	"io"
	glog "log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Logger provides interface for logging library.
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

const readWritePerm = 666

// Config holds configuration for logger.
type Config struct {
	Dir   string `json:"dir" validate:"required"`
	Level int    `json:"level" validate:"number"`
}

// Log hold logger object.
type Log struct {
	Logger *zerolog.Logger
}

func (l *Log) Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	l.Logger.Error().Stack().Err(errors.WithStack(err)).Msgf(format, args...)
}

func (l *Log) Error(ctx context.Context, err error, msg string) {
	l.Logger.Error().Stack().Err(errors.WithStack(err)).Msg(msg)
}

func (l *Log) Fatalf(ctx context.Context, err error, format string, args ...interface{}) {
	l.Logger.Fatal().Stack().Err(errors.WithStack(err)).Msgf(format, args...)
}

func (l *Log) Fatal(ctx context.Context, err error, msg string) {
	l.Logger.Fatal().Stack().Err(errors.WithStack(err)).Msg(msg)
}

func (l *Log) Infof(ctx context.Context, format string, args ...interface{}) {
	l.Logger.Info().Msgf(format, args...)
}

func (l *Log) Info(ctx context.Context, msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *Log) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.Logger.Debug().Msgf(format, args...)
}

func (l *Log) Debug(ctx context.Context, msg string) {
	l.Logger.Debug().Msg(msg)
}

type levelWriter struct {
	app io.WriteCloser
	err io.WriteCloser
}

func (lw *levelWriter) WriteLevel(l zerolog.Level, msg []byte) (n int, err error) {
	w := lw.app
	if l > zerolog.InfoLevel {
		w = lw.err
	}

	n, err = w.Write(msg)

	return n, errors.Wrap(err, "LevelWrite write")
}

func (lw *levelWriter) Write(msg []byte) (n int, err error) {
	n, err = lw.app.Write(msg)

	return n, errors.Wrap(err, "LevelWrite Default Write")
}

func (lw *levelWriter) Close() error {
	if err := lw.app.Close(); err != nil {
		return errors.Wrap(err, "LevelWriter Default Close")
	}

	return errors.Wrap(lw.err.Close(), "LevelWriter Error Close")
}

// NewZeroLogger returns a new instance of zerologger.
func NewZeroLogger(conf Config) (*Log, error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	lvlWriter, err := newLevelWriter(conf.Dir+"/app.log", conf.Dir+"/error.log")
	if err != nil {
		return nil, errors.Wrap(err, "New LevelWriter")
	}

	logWriter := zerolog.New(lvlWriter).Level(zerolog.Level(conf.Level))
	logWriter = logWriter.With().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).Timestamp().Logger()

	return &Log{&logWriter}, nil
}

func newLevelWriter(appFile, errFile string) (*levelWriter, error) {
	defaultWriter, err := os.OpenFile(filepath.Clean(appFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, readWritePerm)
	if err != nil {
		return nil, errors.Wrap(err, "App Log File Open")
	}

	errorWriter, err := os.OpenFile(filepath.Clean(errFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, readWritePerm)

	return &levelWriter{defaultWriter, errorWriter}, errors.Wrap(err, "Error Log File Open")
}

type errorWriter struct {
	logger Logger
}

func (ew *errorWriter) Write(p []byte) (int, error) {
	ew.logger.Error(context.Background(), errors.New(string(p)), "Server Internal")

	return len(p), nil
}

// NewServerLogger return new logger for rest server errors.
func NewServerLogger(logger Logger) *glog.Logger {
	return glog.New(&errorWriter{logger}, "" /* prefix */, 0 /* flags */)
}

func NewConsoleLogger() *Log {
	consoleLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	consoleLogger = consoleLogger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return &Log{&consoleLogger}
}
