package logger

import (
	"context"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type log struct {
	Logger *zerolog.Logger
}

func (l *log) Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	l.Logger.Error().Stack().Err(err).Msgf(format, args...)
}

func (l *log) Error(ctx context.Context, err error, msg string) {
	l.Logger.Error().Stack().Err(err).Msgf(msg)
}

func (l *log) Fatalf(ctx context.Context, err error, format string, args ...interface{}) {
	l.Logger.Fatal().Stack().Err(err).Msgf(format, args...)
}

func (l *log) Fatal(ctx context.Context, err error, msg string) {
	l.Logger.Error().Stack().Err(err).Msg(msg)
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

func NewZeroLogger(conf LogConfig) (Logger, error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.Level(conf.Level))

	defaultWriter, err := os.OpenFile(conf.Dir+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "App Log File Open")
	}

	errorWriter, err := os.OpenFile(conf.Dir+"/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "Error Log File Open")
	}

	levelWriter := LevelWriter{defaultWriter, errorWriter}

	logWriter := zerolog.New(&levelWriter).With().Caller().Timestamp().Logger()

	return &log{&logWriter}, nil
}
