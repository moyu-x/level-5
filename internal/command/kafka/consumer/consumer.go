package consumer

import (
	"context"
	"os"

	"github.com/bytedance/sonic"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/moyu-x/level-5/pkg/config"
	lk "github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/logger"
)

type Config struct {
	Topic     string
	GroupID   string
	ServerAdd string
	Filter    string
	Mode      string
	FilePath  string
}

func Run(configPath string, cc Config) {
	c := config.NewConfig(configPath)
	logger.NewLogger(c)
	k := lk.NewKafka(c)
	r := k.Reader(cc.Topic, cc.GroupID, cc.ServerAdd)

	switch cc.Mode {
	case "console":
		consumer(r, cc.Filter, func(it string) {
			log.Info().Msg(it)
		})
	case "file":
		createFile(cc.FilePath)
		file, _ := os.OpenFile(cc.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		defer file.Close()

		consumer(r, cc.Filter, func(it string) {
			file.WriteString(it + "\n")
		})
	default:
		log.Fatal().Msg("mode is not support")

	}
}

func createFile(filepath string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file, err := os.Create(filepath)
		if err != nil {
			log.Fatal().Msgf("create file has oucur error. reason: %v", err)
		}
		file.Close()
		log.Info().Msgf("create file %s success", filepath)
		return
	}
	log.Info().Msgf("file %s already exists", filepath)
}

func consumer(r *kafka.Reader, filter string, fn func(it string)) {
	var compile *vm.Program
	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Error().Msgf("read http kafka data has oucur error. reason: %v", err)
			continue
		}
		if compile == nil && filter != "" {
			compile, err = expr.Compile(filter)
			if err != nil {
				log.Fatal().Msgf("compile filter has oucur error. reason: %v", err)
			}
		} else if compile != nil && filter != "" {
			var data map[string]interface{}
			err = sonic.UnmarshalString(string(message.Value), &data)
			if err != nil {
				log.Error().Msgf("unmarshal kafka data has oucur error. reason: %v", err)
				continue
			}
			out, err := expr.Run(compile, data)
			if err != nil {
				log.Error().Msgf("run filter has oucur error. reason: %v", err)
				continue
			}
			if !out.(bool) {
				continue
			}
		}

		fn(string(message.Value))
	}
}
