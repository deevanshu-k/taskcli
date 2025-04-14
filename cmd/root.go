package cmd

import (
	"fmt"
	"taskcli/config"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "taskcli",
	Short: "Simple CLI for managing tasks",
	Run: func(cmd *cobra.Command, args []string) {
		config.Load()
		fmt.Println(config.Config)
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
