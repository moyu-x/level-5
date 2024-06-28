package ants

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/panjf2000/ants"
)

func NewAnts(l *logr.Logger) *ants.Pool {
	pool, err := ants.NewPool(128, ants.WithPreAlloc(true))
	if err != nil {
		l.Error(err, "init ants pool error")
		os.Exit(-1)
	}

	return pool
}
