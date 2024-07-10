package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moyu-x/level-5/pkg/config"
)

var kafkaConfig config.KafkaConfig

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "comsumer, producer and manager kafka",
}

func init() {
	rootCmd.AddCommand(kafkaCmd)
}
