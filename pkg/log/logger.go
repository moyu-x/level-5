package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/moyu-x/level-5/pkg/config"
)

func NewLogger(c *config.Bootstrap) *zerolog.Logger {
	switch c.Logger.
		Level {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &logger
}
