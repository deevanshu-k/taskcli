package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	updateCommand.Flags().StringP("description", "d", "", "Description of the task")
	updateCommand.Flags().StringP("status", "s", "", "Status of the task (P/I/C)")

	rootCommand.AddCommand(updateCommand)
}

var updateCommand = &cobra.Command{
	Use:     "update",
	Short:   "Update existing tasks",
	Example: "taskcli update <taskid> -d 'New task desc' -s 'P/I/C' \n P=Pending I=InProgress C=Completed \n -d for updateing description \n -s for updating status",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieving flags
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		// Retrieving the id argument
		id := args[0]

		// Output the task info (this is just for demonstration)
		fmt.Printf("Task ID: %s\n", id)
		fmt.Printf("Description: %s\n", description)
		fmt.Printf("Status: %s\n", status)
	},
}
