package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/deevanshu-k/taskcli/libs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addTaskCommand)
}

var addTaskCommand = &cobra.Command{
	Use:   "add 'task'",
	Short: "Add task to the list.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("pass tasks as arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := libs.CreateTask(args, libs.Pending)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	DisableFlagsInUseLine: true,
}
