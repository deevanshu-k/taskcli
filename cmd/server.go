package cmd

import (
	"taskcli/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(serverCommand)
}

var serverCommand = &cobra.Command{
	Use:   "start",
	Args:  cobra.NoArgs,
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Start()
	},
}
