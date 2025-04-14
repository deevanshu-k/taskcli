package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	listCommand.Flags().StringP("filter", "f", "", "Filter tasks by status")

	rootCommand.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:     "list",
	Short:   "List all tasks",
	Example: "taskcli list",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")
		if filter != "" {
			// Implement filtering logic here
			fmt.Println("List filtered tasks with status:", filter)
		} else {
			// List all tasks
			fmt.Println("List all tasks.")
		}
	},
}
