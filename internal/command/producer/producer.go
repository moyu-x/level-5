package producer

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"

	"github.com/moyu-x/level-5/internal/fake"
	"github.com/moyu-x/level-5/pkg/config"
	lkafka "github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/log"
)

type ProduceConfig struct {
	FilePath   string
	Round      int
	Mode       string
	Data       string
	Topic      string
	FakeType   string
	ServerAddr string
}

func Run(configPath string, p ProduceConfig) {
	c := config.NewConfig(configPath)
	l := log.NewLogger(c)
	k := lkafka.NewKafka(c, l)
	w := k.Writer(p.Topic, p.ServerAddr)

	producer := NewProducer(p, l, w)

	switch p.Mode {
	case "d":
		producer.replayData()
	case "f":
		producer.fakeData()
	default:
		l.Error().Msg("can't found any match mode")
	}
}

type Producer struct {
	pc     ProduceConfig
	logger *zerolog.Logger
	log    *zerolog.Logger
	writer *kafka.Writer
}

func NewProducer(pc ProduceConfig, l *zerolog.Logger, w *kafka.Writer) *Producer {
	return &Producer{pc, l, l, w}
}

func (p *Producer) replayData() {
	if p.pc.Round <= 1000 {
		msgs := messages(p.pc.Round, p.pc.Data)
		err := p.writer.WriteMessages(context.Background(), msgs...)
		if err != nil {
			p.log.Error().Msgf("send kafka data has error. reason: %v", err)
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

func (p *Producer) fakeData() {
	ctx := context.Background()
	p.log.Info().Msg("start to send fake data")
	f := fake.New()
	var msgs []kafka.Message
	count := 0

	ts := time.Now().UnixMilli()

	for i := 0; i < p.pc.Round; i++ {
		if count == 1000 {
			err := p.writer.WriteMessages(ctx, msgs...)
			if err != nil {
				p.log.Error().Msgf("send kafka data has error. reason: %v", err)
			}
			msgs = []kafka.Message{}
			count = 0
		}

		data := p.pc.Data

		data = strings.ReplaceAll(data, "{{sip}}", f.IPv4Address())
		data = strings.ReplaceAll(data, "{{dip}}", f.IPv4Address())
		data = strings.ReplaceAll(data, "{{sport}}", strconv.Itoa(f.Number(0, 65535)))
		data = strings.ReplaceAll(data, "{{dport}}", strconv.Itoa(f.Number(0, 65535)))
		data = strings.ReplaceAll(data, "{{ts}}", strconv.FormatInt(ts, 10))

		fakeData := kafka.Message{
			Value: []byte(data),
		}
		msgs = append(msgs, fakeData)
		count++
	}

	if len(msgs) > 0 {
		p.flushMessage(ctx, msgs)
	}
}

func (p *Producer) flushMessage(ctx context.Context, msgs []kafka.Message) {
	err := p.writer.WriteMessages(ctx, msgs...)
	if err != nil {
		p.log.Error().Msgf("send kafka data has error. reason: %v", err)
	}
}
