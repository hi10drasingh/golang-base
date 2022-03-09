package drmlog

import (
	"context"
	"io"
	glog "log"
	"os"
	"path/filepath"

	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Logger provides interface for logging library
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

// Config holds init dependencies for New Logger
type Config struct {
	LogConfig config.LogConfig
}

type log struct {
	Logger *zerolog.Logger
}

func (l *log) Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	if err != nil {
		l.Logger.Error().Stack().Err(errors.WithStack(err)).Msgf(format, args...)
		return
	}
	l.Logger.Error().Stack().Msgf(format, args...)
}

func (l *log) Error(ctx context.Context, err error, msg string) {
	if err != nil {
		l.Logger.Error().Stack().Err(errors.WithStack(err)).Msg(msg)
		return
	}
	l.Logger.Error().Stack().Msg(msg)
}

func (l *log) Fatalf(ctx context.Context, err error, format string, args ...interface{}) {
	if err != nil {
		l.Logger.Fatal().Stack().Err(errors.WithStack(err)).Msgf(format, args...)
		return
	}
	l.Logger.Fatal().Stack().Msgf(format, args...)
}

func (l *log) Fatal(ctx context.Context, err error, msg string) {
	if err != nil {
		l.Logger.Fatal().Stack().Msg(msg)
		return
	}
	l.Logger.Fatal().Stack().Err(errors.WithStack(err)).Msg(msg)
}

func (l *log) Infof(ctx context.Context, format string, args ...interface{}) {
	l.Logger.Info().Msgf(format, args...)
}

func (l *log) Info(ctx context.Context, msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *log) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.Logger.Debug().Msgf(format, args...)
}

func (l *log) Debug(ctx context.Context, msg string) {
	l.Logger.Debug().Msg(msg)
}

type levelWriter struct {
	io.Writer
	ErrorWriter io.Writer
}

func (lw *levelWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	w := lw.Writer
	if l > zerolog.InfoLevel {
		w = lw.ErrorWriter
	}
	return w.Write(p)
}

// NewZeroLogger returns a new instance of zerologger
func NewZeroLogger(conf Config) (Logger, error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	lvlWriter, err := newLevelWriter(conf.LogConfig.Dir+"/app.log", conf.LogConfig.Dir+"/error.log")
	if err != nil {
		return nil, errors.Wrap(err, "New Level Writer")
	}

	logWriter := zerolog.New(lvlWriter).Level(zerolog.Level(conf.LogConfig.Level)).With().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).Timestamp().Logger()

	return &log{&logWriter}, nil
}

func newLevelWriter(appFile string, errFile string) (*levelWriter, error) {
	defaultWriter, err := os.OpenFile(filepath.Clean(appFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "App Log File Open")
	}

	errorWriter, err := os.OpenFile(filepath.Clean(errFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "Error Log File Open")
	}

	return &levelWriter{defaultWriter, errorWriter}, nil
}

type errorWriter struct {
	logger Logger
}

func (ew *errorWriter) Write(p []byte) (int, error) {
	ew.logger.Error(context.TODO(), errors.New(string(p)), "Server Internal Error")
	return len(p), nil
}

// NewServerLogger return new logger for rest server errors
func NewServerLogger(logger Logger) *glog.Logger {
	return glog.New(&errorWriter{logger}, "" /* prefix */, 0 /* flags */)
}
