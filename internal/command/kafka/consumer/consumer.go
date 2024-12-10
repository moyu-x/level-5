package consumer

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/expr-lang/expr"
	"github.com/rs/zerolog/log"

	"github.com/moyu-x/level-5/pkg/config"
	"github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/logger"
	"github.com/moyu-x/level-5/pkg/pool"
)

type Config struct {
	Topic     string
	GroupID   string
	ServerAdd string
	Filter    string
}

func Run(configPath string, cc Config) {
	c := config.NewConfig(configPath)
	logger.NewLogger(c)
	k := kafka.NewKafka(c)
	r := k.Reader(cc.Topic, cc.GroupID, cc.ServerAdd)
	antsPool := pool.NewAnts()

	if cc.Filter == "" {
		for {
			message, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Error().Msgf("read http kafka data has oucur error. reason: %v", err)
				continue
			}
			log.Info().Msg(string(message.Value))
		}
	} else {
		compile, err := expr.Compile(cc.Filter)
		if err != nil {
			panic(err)
		}

		for {
			message, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Error().Msgf("read http kafka data has oucur error. reason: %v", err)
				continue
			}
			content := string(message.Value)
			_ = antsPool.Submit(func() {
				var data map[string]interface{}
				err = sonic.UnmarshalString(content, &data)
				if err != nil {
					log.Error().Msgf("unmarshal http kafka data has oucur error. reason: %v", err)
					return
				}
				out, err := expr.Run(compile, data)
				if err != nil {
					log.Error().Msgf("run http kafka data has oucur error. reason: %v", err)
				}
				if out.(bool) {
					log.Info().Msg(content)
				}
			})
		}
	}
}
