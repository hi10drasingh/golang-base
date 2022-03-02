package logger

import (
	"io"
	"os"

	"github.com/droomlab/drm-coupon/pkg/config"
	"github.com/rs/zerolog"
)

type log struct {
	Logger *zerolog.Logger
}

func (l *log) Errorf(err error, msg string, args ...interface{}) {
	l.Logger.Error().Stack().Err(err).Msgf(msg, args...)
}

func (l *log) Error(err error, msg string) {
	l.Logger.Error().Stack().Err(err).Msgf(msg)
}

func (l *log) Fatalf(err error, msg string, args ...interface{}) {
	l.Logger.Fatal().Stack().Err(err).Msgf(msg, args...)
}

func (l *log) Fatal(err error, msg string) {
	l.Logger.Error().Stack().Err(err).Msg(msg)
}

func (l *log) Infof(msg string, args ...interface{}) {
	l.Logger.Info().Msgf(msg, args...)
}

func (l *log) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *log) Debugf(msg string, args ...interface{}) {
	l.Logger.Debug().Msgf(msg, args...)
}

func (l *log) Debug(msg string) {
	l.Logger.Debug().Msgf(msg)
}

type LevelWriter struct {
	io.Writer
	ErrorWriter io.Writer
}

func (lw *LevelWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	w := lw.Writer
	if l > zerolog.InfoLevel {
		w = lw.ErrorWriter
	}
	return w.Write(p)
}

func NewZeroLogger(conf config.Logging) (Logger, error) {
	defaultWriter, err := os.OpenFile(conf.Dir+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	errorWriter, err := os.OpenFile(conf.Dir+"/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	levelWriter := LevelWriter{defaultWriter, errorWriter}

	logWriter := zerolog.New(levelWriter).With().Timestamp().Logger()

	return &log{&logWriter}, nil
}
