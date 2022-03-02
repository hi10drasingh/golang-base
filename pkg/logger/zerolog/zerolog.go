package zerolog

import (
	"io"
	"os"

	"github.com/droomlab/drm-coupon/pkg/config"
	"github.com/droomlab/drm-coupon/pkg/logger"
	"github.com/rs/zerolog"
)

type Log struct {
	Logger *zerolog.Logger
}

func (l *Log) Errorf(err error, msg string, args ...interface{}) {
	l.Logger.Error().Stack().Err(err).Msgf(msg, args)
}

func (l *Log) Error(err error, msg string) {
	l.Logger.Error().Stack().Err(err).Msgf(msg)
}

func (l *Log) Fatalf(err error, msg string, args ...interface{}) {
	l.Logger.Fatal().Stack().Err(err).Msgf(msg, args)
}

func (l *Log) Fatal(err error, msg string) {
	l.Logger.Error().Stack().Err(err).Msg(msg)
}

func (l *Log) Infof(msg string, args ...interface{}) {
	l.Logger.Info().Msgf(msg, args)
}

func (l *Log) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *Log) Debugf(msg string, args ...interface{}) {
	l.Logger.Debug().Msgf(msg, args)
}

func (l *Log) Debug(msg string) {
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

func NewZeroLogger(lc *config.Logging) (*logger.Logger, error) {
	defaultLogger, err := getDefaultLogger(lc)
	if err != nil {
		return nil, err
	}

	errorLogger, err := getErrorLogger(lc)
	if err != nil {
		return nil, err
	}
	
	levelWriter := LevelWriter{defaultLogger, errorLogger}

	writer := zerolog.New(levelWriter).With().Timestamp().Logger()


	return &Log{writer}, nil
}

func getDefaultLogger (lc *config.Logging) (io.Writer, error) {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(lc.Dir+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file, nil
}

func getErrorLogger (lc *config.Logging) (io.Writer, error) {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(lc.Dir+"/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file, nil
}