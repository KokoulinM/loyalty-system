package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Log(msg string)
	Fatal(msg string)
	Error(msg string)
}

type logger struct {
	zl zerolog.Logger
}

func New(level zerolog.Level) *logger {
	return &logger{
		zl: zerolog.New(os.Stdout).Level(level),
	}
}

func (l *logger) Log(msg string) {
	l.zl.Log().Msg(msg)
}

func (l *logger) Fatal(msg string) {
	l.zl.Fatal().Msg(msg)
}

func (l *logger) Error(msg string) {
	l.zl.Error().Msg(msg)
}
