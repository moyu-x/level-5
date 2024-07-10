/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moyu-x/level-5/internal/command/consumer"
)

var cc consumer.ConsumerConfig

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:     "consumer",
	Aliases: []string{"c"},
	Short:   "consumer kafka message ",
	Run: func(cmd *cobra.Command, args []string) {
		consumer.Run(configPath, cc)
	},
}

func init() {
	kafkaCmd.AddCommand(consumerCmd)
	consumerCmd.Flags().StringVarP(&cc.Topic, "topic", "t", "", "topic name")
	consumerCmd.Flags().StringVarP(&cc.GroupID, "group-id", "g", "", "group id")
}
