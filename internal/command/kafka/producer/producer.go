package producer

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gookit/goutil/maputil"
	"github.com/panjf2000/ants"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/moyu-x/level-5/internal/fake"
	"github.com/moyu-x/level-5/pkg/config"
	lkafka "github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/logger"
	"github.com/moyu-x/level-5/pkg/pool"
)

type Config struct {
	FilePath   string
	Round      int
	Mode       string
	Data       string
	Topic      string
	FakeType   string
	ServerAddr string
	BatchSize  int
}

func Run(configPath string, p Config) {
	ctx := context.Background()
	c := config.NewConfig(configPath)
	k := lkafka.NewKafka(c)
	antsPool := pool.NewAnts()
	logger.NewLogger(c)
	w := k.Writer(p.Topic, p.ServerAddr)

	producer := NewProducer(ctx, p, w, antsPool)

	switch p.Mode {
	case "r":
		producer.replayData(ctx)
	case "f":
		producer.fakeData()
	case "i":
		producer.fromFile()
	default:
		log.Error().Msg("can't found any match mode")
	}
}

type Producer struct {
	ctx    context.Context
	pc     Config
	writer *kafka.Writer
	pool   *ants.Pool
}

func NewProducer(ctx context.Context, pc Config, w *kafka.Writer, ants *ants.Pool) *Producer {
	return &Producer{ctx, pc, w, ants}
}

func (p *Producer) replayData(ctx context.Context) {
	round := p.pc.Round / p.pc.BatchSize
	msgs := messages(p.pc.BatchSize, p.pc.Data)
	for i := range round {
		err := p.writer.WriteMessages(ctx, msgs...)
		if err != nil {
			log.Error().Msgf("send kafka data has error. reason: %v", err)
		}
		log.Info().Msgf("round %d, cap %d", i, len(msgs))
	}

	least := p.pc.Round - p.pc.BatchSize*round
	if least > 0 {
		err := p.writer.WriteMessages(ctx, msgs[:least]...)
		if err != nil {
			log.Error().Msgf("send kafka data has error. reason: %v", err)
		}
	}
}

func messages(size int, data string) []kafka.Message {
	var msgs []kafka.Message
	for range size {
		m := kafka.Message{
			Value: []byte(data),
		}
		msgs = append(msgs, m)
	}
	return msgs
}

func (p *Producer) fakeData() {
	ctx := context.Background()
	log.Info().Msg("start to send fake data")
	f := fake.New()
	var msgs []kafka.Message
	count := 0

	ts := time.Now().UnixMilli()

	for range p.pc.Round {
		if count == 1000 {
			err := p.writer.WriteMessages(ctx, msgs...)
			if err != nil {
				log.Error().Msgf("send kafka data has error. reason: %v", err)
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
	log.Info().Msgf("flush %d messages", len(msgs))
	err := p.writer.WriteMessages(ctx, msgs...)
	if err != nil {
		log.Error().Msgf("send kafka data has error. reason: %v", err)
	}
}

func (p *Producer) fromFile() {
	// read file from one by one
	log.Info().Msg("start to send data from file")
	file, err := os.Open(p.pc.FilePath)
	if err != nil {
		log.Fatal().Msg("can't open file")
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var msgs []kafka.Message
	count := 0
	for scanner.Scan() {
		data := scanner.Text()
		var raw map[string]any
		err := sonic.UnmarshalString(data, &raw)
		if err != nil {
			log.Error().Msgf("can't unmarshal data. reason: %v", err)
			continue
		}
		if maputil.HasKey(raw, "ts") {
			raw["ts"] = time.Now().UnixMilli()
		}

		data, err = sonic.MarshalString(raw)
		if err != nil {
			log.Error().Msgf("can't marshal data. reason: %v", err)
			continue
		}

		msgs = append(msgs, kafka.Message{
			Value: []byte(data),
		})
		count++

		if count >= p.pc.BatchSize {
			p.flushMessage(p.ctx, msgs)
			msgs = []kafka.Message{}
			count = 0
		}
	}
	p.flushMessage(p.ctx, msgs)
}
