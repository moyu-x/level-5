package ants

import (
	"os"
	"runtime"

	"github.com/panjf2000/ants"
	"github.com/rs/zerolog"
)

func NewAnts(l *zerolog.Logger) *ants.Pool {
	pool, err := ants.NewPool(runtime.NumCPU()*2, ants.WithPreAlloc(true))
	if err != nil {
		l.Error().Msgf("init ants pool error. reason: %v", err)
		os.Exit(-1)
	}

	return pool
}
