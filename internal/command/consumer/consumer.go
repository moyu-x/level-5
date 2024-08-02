package consumer

import (
	"context"

	"github.com/moyu-x/level-5/pkg/config"
	"github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/log"
)

type Config struct {
	Topic     string
	GroupID   string
	ServerAdd string
}

func Run(configPath string, cc Config) {
	c := config.NewConfig(configPath)
	l := log.NewLogger(c)
	k := kafka.NewKafka(c, l)
	r := k.Reader(cc.Topic, cc.GroupID, cc.ServerAdd)
	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			l.Error().Msgf("read http kafka data has oucur error. reason: %v", err)
			continue
		}
		l.Info().Msg(string(message.Value))
	}
}
