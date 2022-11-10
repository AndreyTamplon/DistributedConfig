package logger

import (
	"github.com/rs/zerolog"
	"os"
	"strings"
)

type Interface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()
	zerolog.SetGlobalLevel(l)
	return &Logger{
		logger: &logger,
	}
}

func (logger *Logger) Debug(message string, args ...interface{}) {
	logger.logger.Debug().Msgf(message, args...)
}

func (logger *Logger) Info(message string, args ...interface{}) {
	logger.logger.Info().Msgf(message, args...)
}

func (logger *Logger) Warn(message string, args ...interface{}) {
	logger.logger.Warn().Msgf(message, args...)
}

func (logger *Logger) Error(message string, args ...interface{}) {
	logger.logger.Error().Msgf(message, args...)
}

func (logger *Logger) Fatal(message string, args ...interface{}) {
	logger.logger.Fatal().Msgf(message, args...)
	os.Exit(1)
}
