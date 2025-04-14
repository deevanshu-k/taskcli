package cmd

import (
	"fmt"
	"taskcli/structs"

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

		for _, task := range args {
			command := structs.NewCommand(structs.ADD, nil, &task, nil, nil)

			_, err := command.SendCommand()
			if err != nil {
				fmt.Printf("%v", err)
				return
			}

		}
	},
}
