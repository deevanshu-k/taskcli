package cmd

import (
	"fmt"
	"strconv"
	"taskcli/structs"

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
		if len(args) == 0 {
			cmd.Help()
			return
		}
		taskId, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}

		if description == "" && (len(status) == 0 || len(status) > 1 || status[0] != 'P' && status[0] != 'I' && status[0] != 'C') {
			cmd.Help()
			return
		}

		var rstatus *structs.Status = nil

		if len(status) == 1 && rune(status[0]) == rune(structs.PENDING) {
			temp := structs.PENDING
			rstatus = &temp
		}
		if len(status) == 1 && rune(status[0]) == rune(structs.INPROGRESS) {
			temp := structs.INPROGRESS
			rstatus = &temp
		}
		if len(status) == 1 && rune(status[0]) == rune(structs.COMPLETED) {
			temp := structs.COMPLETED
			rstatus = &temp
		}

		var pdescription *string = nil
		if description != "" {
			pdescription = &description
		}

		command := structs.NewCommand(structs.UPDATE, &taskId, pdescription, rstatus, nil)

		res, err := command.SendCommand()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Println(res)
	},
}
