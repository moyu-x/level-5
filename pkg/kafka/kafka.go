package kafka

import (
	"time"

	"github.com/go-logr/logr"
	kafka "github.com/segmentio/kafka-go"

	"github.com/moyu-x/level-5/pkg/config"
)

type K struct {
	c *config.Bootstrap
	l *logr.Logger
}

func NewKafka(c *config.Bootstrap, l *logr.Logger) *K {
	kafka := &K{
		c: c, l: l,
	}

	return kafka
}

func (k *K) dialer() *kafka.Dialer {
	return &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       Tls(k.l, k.c),
	}
}

func (k *K) Reader() *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               []string{k.c.Kafka.ServerAddr},
		GroupID:               k.c.Kafka.GroupId,
		GroupTopics:           []string{k.c.Kafka.Topic},
		MaxBytes:              10e6,
		CommitInterval:        time.Second,
		Dialer:                k.dialer(),
		QueueCapacity:         1024,
		WatchPartitionChanges: true,
		StartOffset:           kafka.LastOffset,
	})
	return reader
}

func (k *K) Writer() *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          []string{k.c.Kafka.ServerAddr},
		Topic:            k.c.Kafka.Topic,
		Balancer:         &kafka.Hash{},
		Dialer:           k.dialer(),
		CompressionCodec: kafka.Zstd.Codec(),
	})
	return w
}
