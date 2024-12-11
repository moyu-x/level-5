package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moyu-x/level-5/internal/command/kafka/admin"
)

var adminConfig admin.Config

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "kafka admin command",
	Run: func(cmd *cobra.Command, args []string) {
		admin.Run(configPath, adminConfig)
	},
}

func init() {
	kafkaCmd.AddCommand(adminCmd)
	adminCmd.Flags().StringVarP(&adminConfig.ServerAddr, "server-add", "s", "", "server address")
	adminCmd.Flags().StringVarP(&adminConfig.Mode, "mode", "m", "list-topic", "admin mode")
	adminCmd.Flags().IntVarP(&adminConfig.Partition, "partition", "p", 9, "partition")
	adminCmd.Flags().IntVarP(&adminConfig.Replication, "replication", "r", 1, "replication")
	adminCmd.Flags().StringVarP(&adminConfig.Topic, "topic", "t", "", "topic name")
	adminCmd.Flags().StringVarP(&adminConfig.GroupId, "group-id", "g", "", "group id")
}
