package kafka

import (
	"context"
	"os"
	"time"

	"github.com/gookit/goutil/maputil"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/moyu-x/level-5/pkg/common"
	"github.com/moyu-x/level-5/pkg/config"
)

type K struct {
	c *config.Bootstrap
}

func NewKafka(c *config.Bootstrap) *K {
	k := &K{
		c: c,
	}

	return k
}

func (k *K) dialer() *kafka.Dialer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       nil,
	}
	if k.c.Kafka != nil && k.c.Kafka.KeyFilePath != "" && k.c.Kafka.CertFilePath != "" {
		dialer.TLS = Tls(k.c)
	}

	return dialer
}

func (k *K) dial(serverAddr string) *kafka.Conn {
	conn, err := kafka.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal().Msgf("failed to dial Kafka server: %v", err)
	}

	return conn
}

func (k *K) Reader(topic string, groupId string, serverAddr string) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               []string{serverAddr},
		GroupID:               groupId,
		GroupTopics:           []string{topic},
		MaxBytes:              10e6,
		CommitInterval:        time.Second,
		Dialer:                k.dialer(),
		QueueCapacity:         1024,
		WatchPartitionChanges: true,
		StartOffset:           kafka.LastOffset,
		RetentionTime:         time.Minute * 1,
		Logger:                kafka.LoggerFunc(k.infoF),
		ErrorLogger:           kafka.LoggerFunc(k.errorF),
	})
	return reader
}

func (k *K) Writer(topic string, serverAddr string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(serverAddr),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		Compression:  kafka.Zstd,
		Logger:       kafka.LoggerFunc(k.infoF),
		ErrorLogger:  kafka.LoggerFunc(k.errorF),
		RequiredAcks: 0,
	}

	if k.c.Kafka != nil && k.c.Kafka.KeyFilePath != "" && k.c.Kafka.CertFilePath != "" {
		w.Transport = &kafka.Transport{
			TLS: Tls(k.c),
		}
	}
	return w
}

func (k *K) ListTopics(serverAddr string) {
	dial := k.dialer()
	conn, err := dial.DialContext(context.Background(), "tcp", serverAddr)
	if err != nil {
		log.Fatal().Msgf("failed to dial Kafka server: %v", err)
	}
	defer conn.Close()
	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]int{}

	for _, p := range partitions {
		if p.Topic == "__consumer_offsets" {
			continue
		}
		topicName := p.Topic
		if maputil.HasKey(m, topicName) {
			m[topicName] = m[topicName] + 1
			continue
		}
		m[p.Topic] = 1
	}

	k.renderTable(m)
}

func (k *K) renderTable(m map[string]int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"topic name", "partition count"})

	for i := range m {
		t.AppendRow(table.Row{i, m[i]})
	}
	t.Render()
}

func (k *K) infoF(msg string, a ...interface{}) {
	if k.c.Logger.KafkaLevel == common.INFO {
		log.Info().Msgf(msg, a...)
	}
}

func (k *K) errorF(msg string, a ...interface{}) {
	log.Error().Msgf(msg, a...)
}
