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

func (k *K) conn(serverAddr string) *kafka.Conn {
	dial := k.dialer()
	conn, err := dial.DialContext(context.Background(), "tcp", serverAddr)
	if err != nil {
		log.Fatal().Msgf("failed to dial Kafka server: %v", err)
	}

	return conn
}

func (k *K) client(serverAddr string) *kafka.Client {
	client := kafka.Client{
		Addr:    kafka.TCP(serverAddr),
		Timeout: 10 * time.Second,
	}
	if k.c.Kafka != nil && k.c.Kafka.KeyFilePath != "" && k.c.Kafka.CertFilePath != "" {
		client.Transport = &kafka.Transport{
			TLS: Tls(k.c),
		}
	}
	return &client
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
	topics := k.listTopics(serverAddr)
	k.renderTable(topics)
}

func (k *K) listTopics(serverAddr string) map[string]int {
	conn := k.conn(serverAddr)
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

	return m
}

// CreateTopic 创建 topic
func (k *K) CreateTopic(topic string, partition int, replication int, serverAddr string) {
	conn := k.conn(serverAddr)
	defer conn.Close()

	err := conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partition,
		ReplicationFactor: replication,
	})
	if err != nil {
		log.Fatal().Msgf("create topic failed: %v", err)
	}
}

// DeleteTopic 删除 topic
func (k *K) DeleteTopic(topic string, serverAddr string) {
	conn := k.conn(serverAddr)
	defer conn.Close()

	err := conn.DeleteTopics(topic)
	if err != nil {
		log.Fatal().Msgf("delete topic failed: %v", err)
	}
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

func (k *K) ListConsumerGroups(serverAddr string) {
	dial := k.dialer()
	conn, err := dial.DialContext(context.Background(), "tcp", serverAddr)
	if err != nil {
		log.Fatal().Msgf("failed to dial Kafka server: %v", err)
	}
	defer conn.Close()
	partitions, err := conn.ReadPartitions()
	if err != nil {
		log.Fatal().Msgf("list consumer group failed: %v", err)
	}

	for _, partition := range partitions {
		partConn, err := dial.DialPartition(context.Background(), "tcp", serverAddr, kafka.Partition{
			Topic: partition.Topic,
			ID:    partition.ID,
		})
		if err != nil {
			log.Error().Msgf("list consumer group failed: %v", err)
			continue
		}
		defer partConn.Close()
		offset, err := partConn.ReadLastOffset()
		if err != nil {
			log.Error().Msgf("list consumer group failed: %v", err)
			continue
		}
		log.Info().Msgf("topic: %s, partition: %d, offset: %d", partition.Topic, partition.ID, offset)
	}
}

// DeleteGroup 删除 consumer group
func (k *K) DeleteGroup(groupId string, serverAddr string) {
	client := k.client(serverAddr)
	req := kafka.DeleteGroupsRequest{
		GroupIDs: []string{groupId},
	}
	_, err := client.DeleteGroups(context.Background(), &req)
	if err != nil {
		log.Error().Msgf("delete consumer group failed: %v", err)
	}
}

func (k *K) infoF(msg string, a ...interface{}) {
	if k.c.Logger.KafkaLevel == common.INFO {
		log.Info().Msgf(msg, a...)
	}
}

func (k *K) errorF(msg string, a ...interface{}) {
	log.Error().Msgf(msg, a...)
}
