package log

import (
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"

	"github.com/moyu-x/level-5/pkg/config"
)

func NewLogger(c *config.Bootstrap) *logr.Logger {
	switch c.Logger.
		Level {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	logger := zerolog.New(output).With().Timestamp().Logger()
	log := zerologr.New(&logger)
	return &log
}
