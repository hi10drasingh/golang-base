package logger

import (
	"io"
	"os"

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
	l.Logger.Error().Stack().Err(err).Msgf(format, args...)
}

func (l *log) Error(err error, msg string) {
	l.Logger.Error().Stack().Err(err).Msgf(msg)
}

func (l *log) Fatalf(err error, format string, args ...interface{}) {
	l.Logger.Fatal().Stack().Err(err).Msgf(format, args...)
}

func (l *log) Fatal(err error, msg string) {
	l.Logger.Error().Stack().Err(err).Msg(msg)
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
	l.Logger.Debug().Msgf(msg)
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
	zerolog.SetGlobalLevel(zerolog.Level(conf.Level))

	defaultWriter, err := os.OpenFile(conf.Dir+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "App Log File Open")
	}

	errorWriter, err := os.OpenFile(conf.Dir+"/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "Error Log File Open")
	}

	lvlWriter := levelWriter{defaultWriter, errorWriter}

	logWriter := zerolog.New(&lvlWriter).With().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).Timestamp().Logger()

	return &log{&logWriter}, nil
}
