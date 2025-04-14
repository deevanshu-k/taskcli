package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	deleteCommand.Flags().BoolP("all", "a", false, "Delete all tasks")

	rootCommand.AddCommand(deleteCommand)
}

var deleteCommand = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a task",
	Example: "taskcli delete <task_id>",
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		if all {
			fmt.Println("Delete all tasks.")
		} else {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			taskID := args[0]
			fmt.Println("Delete task with ID:", taskID)
		}
	},
}
