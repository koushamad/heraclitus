package cmd

import (
	"github.com/koushamad/heraclitus/app/server/service"
	"github.com/koushamad/heraclitus/pkg/application"
	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "listen to ttp requests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		service.Listen(port)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", application.Config.GetInt("port"), "port to listen to")
	rootCmd.AddCommand(serveCmd)
}
