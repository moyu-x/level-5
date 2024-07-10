package cmd

import (
	"github.com/spf13/cobra"
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "comsumer, producer and manager kafka",
}

func init() {
	rootCmd.AddCommand(kafkaCmd)
}
