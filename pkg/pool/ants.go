package pool

import (
	"os"
	"runtime"

	"github.com/panjf2000/ants"
	"github.com/rs/zerolog/log"
)

func NewAnts() *ants.Pool {
	pool, err := ants.NewPool(runtime.NumCPU()*2, ants.WithPreAlloc(true))
	if err != nil {
		log.Error().Msgf("init pool pool error. reason: %v", err)
		os.Exit(-1)
	}

	return pool
}
