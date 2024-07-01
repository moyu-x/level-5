package producer

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"

	"github.com/moyu-x/level-5/pkg/config"
	k "github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/log"
)

type ProduceConfig struct {
	FilePath string
	Round    int
	Mode     string
	Data     string
}

func Run(configPath string, p ProduceConfig) {
	c := config.NewConfig(configPath)
	l := log.NewLogger(c)
	k := k.NewKafka(c, l)
	w := k.Writer()

	switch p.Mode {
	case "d":
		replayData(p, w, l)
	default:
		l.Error().Msg("can't found any match mode")
	}
}

func replayData(p ProduceConfig, k *kafka.Writer, l *zerolog.Logger) {
	if p.Round <= 1000 {
		msgs := messages(p.Round, p.Data)
		err := k.WriteMessages(context.Background(), msgs...)
		if err != nil {
			l.Error().Msgf("send kafka data has error. reason: %v", err)
		}
	}
}

func messages(size int, data string) []kafka.Message {
	var msgs []kafka.Message
	for z := 0; z < size; z++ {
		m := kafka.Message{
			Value: []byte(data),
		}
		msgs = append(msgs, m)
	}
	return msgs
}
