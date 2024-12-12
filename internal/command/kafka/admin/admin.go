package admin

import (
	"github.com/rs/zerolog/log"

	"github.com/moyu-x/level-5/pkg/config"
	"github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/logger"
)

type Config struct {
	ServerAddr  string `json:"server_addr"`
	Mode        string `json:"mode"`
	Topic       string `json:"topic"`
	Partition   int    `json:"partition"`
	Replication int    `json:"replication"`
	GroupId     string `json:"group_id"`
}

func Run(globalCfgPath string, cfg Config) {
	c := config.NewConfig(globalCfgPath)
	logger.NewLogger(c)
	k := kafka.NewKafka(c)

	switch cfg.Mode {
	case "list-topic":
		k.ListTopics(cfg.ServerAddr)
	case "delete-topic":
		k.DeleteTopic(cfg.Topic, cfg.ServerAddr)
	case "list-group":
		k.ListConsumerGroups(cfg.ServerAddr)
	case "delete-group":
		k.DeleteGroup(cfg.GroupId, cfg.ServerAddr)
	case "create-topic":
		k.CreateTopic(cfg.Topic, cfg.Partition, cfg.Replication, cfg.ServerAddr)
	default:
		log.Error().Msgf("unknown mode: %s", cfg.Mode)
	}
}
