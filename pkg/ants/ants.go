package ants

import (
	"os"

	"github.com/panjf2000/ants"
	"github.com/rs/zerolog"
)

func NewAnts(l *zerolog.Logger) *ants.Pool {
	pool, err := ants.NewPool(128, ants.WithPreAlloc(true))
	if err != nil {
		l.Error().Msgf("init ants pool error. reason: %v", err)
		os.Exit(-1)
	}

	return pool
}
