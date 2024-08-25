package cmd

import (
	"errors"
	"os"

	"github.com/deevanshu-k/taskcli/libs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addTaskCommand)
}

var addTaskCommand = &cobra.Command{
	Use:   "add",
	Short: "Add task to the list.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("only 1 arguments is allowed")
		}
		if len(args[0]) == 0 {
			return errors.New("task should have minimum 5 charactor")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := libs.CreateTask(args[0], libs.Pending)
		if err != nil {
			println(err)
			os.Exit(1)
		}
	},
}
