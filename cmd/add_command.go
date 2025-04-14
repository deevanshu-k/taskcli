package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(addCommand)
}

var addCommand = &cobra.Command{
	Use:     "add",
	Short:   "Add a new task",
	Example: "taskcli add <task1> <task2>... \n Minumum 1 task is required",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}
