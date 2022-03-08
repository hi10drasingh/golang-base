package logger

import (
	"io"
	glog "log"
	"os"
	"path/filepath"

	"github.com/droomlab/drm-coupon/pkg/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Logger provides interface for logging library
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

type log struct {
	Logger *zerolog.Logger
}

func (l *log) Errorf(err error, format string, args ...interface{}) {
	l.Logger.Error().Stack().Err(errors.WithStack(err)).Msgf(format, args...)
}

func (l *log) Error(err error, msg string) {
	l.Logger.Error().Stack().Err(errors.WithStack(err)).Msg(msg)
}

func (l *log) Fatalf(err error, format string, args ...interface{}) {
	l.Logger.Fatal().Stack().Err(errors.WithStack(err)).Msgf(format, args...)
}

func (l *log) Fatal(err error, msg string) {
	l.Logger.Error().Stack().Err(errors.WithStack(err)).Msg(msg)
}

func (l *log) Infof(format string, args ...interface{}) {
	l.Logger.Info().Msgf(format, args...)
}

func (l *log) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *log) Debugf(format string, args ...interface{}) {
	l.Logger.Debug().Msgf(format, args...)
}

func (l *log) Debug(msg string) {
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
func NewZeroLogger(conf config.LogConfig) (Logger, error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	lvlWriter, err := newLevelWriter(conf.Dir+"/app.log", conf.Dir+"/error.log")
	if err != nil {
		return nil, errors.Wrap(err, "New Level Writer")
	}

	logWriter := zerolog.New(lvlWriter).Level(zerolog.Level(conf.Level)).With().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).Timestamp().Logger()

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
	ew.logger.Error(errors.New(string(p)), "Server Internal Error")
	return len(p), nil
}

// NewServerLogger return new logger for rest server errors
func NewServerLogger(logger Logger) *glog.Logger {
	return glog.New(&errorWriter{logger}, "" /* prefix */, 0 /* flags */)
}
