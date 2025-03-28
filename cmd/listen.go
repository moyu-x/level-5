package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moyu-x/level-5/internal/command/listen"
)

var (
	protocol string
	port     int
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listen a port use tcp/udp",
	Run: func(cmd *cobra.Command, args []string) {
		listen.Run(port, protocol)
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.Flags().StringVarP(&protocol, "protocol", "p", "tcp", "protocol to listen (tcp/udp)")
	listenCmd.Flags().IntVarP(&port, "port", "P", 8080, "port to listen")
}
