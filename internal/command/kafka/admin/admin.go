package admin

import (
	"github.com/moyu-x/level-5/pkg/config"
	"github.com/moyu-x/level-5/pkg/kafka"
	"github.com/moyu-x/level-5/pkg/logger"
)

type Config struct {
	ServerAddr string `json:"server_addr"`
}

func Run(globalCfgPath string, cfg Config) {
	c := config.NewConfig(globalCfgPath)
	logger.NewLogger(c)
	k := kafka.NewKafka(c)

	k.ListTopics(cfg.ServerAddr)
}
