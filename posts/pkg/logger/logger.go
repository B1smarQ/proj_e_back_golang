package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

func New(level string) *Logger {
	return NewLogger(level, nil)
}

func NewLogger(level string, w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(w).With().Timestamp().Logger()

	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug().Msgf(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn().Msgf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Error().Msgf(msg, args...)
	} else {
		l.logger.Error().Msg(msg)
	}
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)
	os.Exit(1)
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.logger.Error().Msg(msg.Error())
	case string:
		l.logger.Info().Msgf(msg, args...)
	default:
		l.logger.Info().Msgf("%s message %v has unknown type %v", level, message, msg)
	}
}
