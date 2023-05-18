package cmd

import (
	"fmt"
	"github.com/koushamad/heraclitus/pkg/application"
	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "listen to ttp requests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called", port)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", application.Config.GetInt("port"), "port to listen to")
	rootCmd.AddCommand(serveCmd)
}
