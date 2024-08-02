package consumer

import (
	"context"
	"strings"

	"github.com/moyu-x/level-5/pkg/ants"
	"github.com/moyu-x/level-5/pkg/config"
	"github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/log"
)

type Config struct {
	Topic     string
	GroupID   string
	ServerAdd string
	Filter    string
}

func Run(configPath string, cc Config) {
	c := config.NewConfig(configPath)
	l := log.NewLogger(c)
	k := kafka.NewKafka(c, l)
	r := k.Reader(cc.Topic, cc.GroupID, cc.ServerAdd)
	pool := ants.NewAnts(l)
	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			l.Error().Msgf("read http kafka data has oucur error. reason: %v", err)
			continue
		}
		content := string(message.Value)
		_ = pool.Submit(func() {
			if cc.Filter == "" {
				l.Info().Msg(content)
			} else {
				if strings.Contains(content, cc.Filter) {
					l.Info().Msg(content)
				}
			}
		})
	}
}
