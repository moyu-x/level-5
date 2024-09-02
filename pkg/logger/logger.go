package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/moyu-x/level-5/pkg/config"
)

func NewLogger(c *config.Bootstrap) {
	switch c.Logger.
		Level {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
