package cmd

import (
	"fmt"
	"taskcli/config"

	"github.com/spf13/cobra"
)

func init() {
	config.Load()
}

var rootCommand = &cobra.Command{
	Use:   "taskcli",
	Short: "Simple CLI for managing tasks",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
	}
}
