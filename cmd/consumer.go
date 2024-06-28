/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moyu-x/level-5/internal/command/consumer"
)

var configPath string

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:     "consumer",
	Aliases: []string{"c"},
	Short:   "consumer kafka message ",
	Run: func(cmd *cobra.Command, args []string) {
		consumer.Run(configPath)
	},
}

func init() {
	kafkaCmd.AddCommand(consumerCmd)
	consumerCmd.Flags().StringVarP(&configPath, "config", "c", "configs/config.toml", "config path")
}
