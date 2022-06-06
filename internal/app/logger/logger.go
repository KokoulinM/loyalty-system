package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level zerolog.Level) *zerolog.Logger {
	log := zerolog.New(os.Stdout).Level(level)

	return &log
}
