package cmd

import (
	"fmt"
	"log"
	"strconv"
	"taskcli/structs"

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
		if !all && len(args) == 0 {
			cmd.Help()
			return
		}

		var taskId *int
		if !all && len(args) > 0 {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalln("Invalid task ID")

			}
			taskId = &id
		}

		command := structs.NewCommand(structs.DELETE, taskId, nil, nil, &all)

		res, err := command.SendCommand()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Println(res)
	},
}
