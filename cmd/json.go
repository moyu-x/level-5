package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moyu-x/level-5/internal/command/json"
)

var jsonConfig json.Config

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "level-5 json tool",
	Run: func(cmd *cobra.Command, args []string) {
		json.Run(jsonConfig)
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().StringVarP(&jsonConfig.InputPath, "input", "i", "", "input path")
	jsonCmd.Flags().StringVarP(&jsonConfig.OutputPath, "output", "o", "", "output path")
}
