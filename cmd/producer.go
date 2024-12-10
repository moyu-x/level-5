package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moyu-x/level-5/internal/command/kafka/producer"
)

var producerConfig producer.Config

// producerCmd represents the producer command
var producerCmd = &cobra.Command{
	Use:     "producer",
	Aliases: []string{"p"},
	Short:   "producer data into kafka",
	Run: func(cmd *cobra.Command, args []string) {
		producer.Run(configPath, producerConfig)
	},
}

func init() {
	kafkaCmd.AddCommand(producerCmd)
	producerCmd.Flags().StringVarP(&producerConfig.FilePath, "path", "p", "", "template or data path")
	producerCmd.Flags().IntVarP(&producerConfig.Round, "round", "r", 1, "replay round")
	producerCmd.Flags().StringVarP(&producerConfig.Mode, "mode", "m", "r", "producer mode: r -> replay, f -> fake data")
	producerCmd.Flags().StringVarP(&producerConfig.Data, "data", "d", "", "kafka data")
	producerCmd.Flags().StringVarP(&producerConfig.Topic, "topic", "t", "", "kafka topic")
	producerCmd.Flags().StringVarP(&producerConfig.FakeType, "type", "f", "", "fake type")
	producerCmd.Flags().StringVarP(&producerConfig.ServerAddr, "server-address", "s", "", "server address")
	producerCmd.Flags().IntVarP(&producerConfig.BatchSize, "batch-size", "b", 1024, "batch size")
}
