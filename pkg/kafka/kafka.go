package kafka

import (
	"time"

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

func (k *K) infoF(msg string, a ...interface{}) {
	if k.c.Logger.KafkaLevel == common.INFO {
		log.Info().Msgf(msg, a...)
	}
}

func (k *K) errorF(msg string, a ...interface{}) {
	log.Error().Msgf(msg, a...)
}
