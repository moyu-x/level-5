package consumer

import (
	"context"

	"github.com/moyu-x/level-5/pkg/config"
	"github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/log"
)

func Run(configPath string) {
	c := config.NewConfig(configPath)
	l := log.NewLogger(c)
	k := kafka.NewKafka(c, l)
	r := k.Reader()
	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			l.Error().Msgf("read http kafka data has oucur error. reason: %v", err)
			continue
		}
		l.Info().Msg(string(message.Value))

	}
}
