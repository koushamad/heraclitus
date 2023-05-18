package cmd

import (
	"github.com/koushamad/heraclitus/app/agent/service"
	"github.com/koushamad/heraclitus/pkg/application"
	"github.com/spf13/cobra"
)

var host string
var port int

var serveCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to the server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		service.ConnectToWebSocket(host, port)
	},
}

func init() {
	serveCmd.Flags().StringVarP(&host, "host", "a", application.Config.GetString("host"), "host to connect to")
	serveCmd.Flags().IntVarP(&port, "port", "p", application.Config.GetInt("port"), "port to connect to")
	rootCmd.AddCommand(serveCmd)
}
