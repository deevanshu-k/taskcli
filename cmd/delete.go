package cmd

import (
	"errors"
	"fmt"

	"github.com/deevanshu-k/taskcli/libs"
	"github.com/spf13/cobra"
)

var deleteAll bool

func init() {
	deleteCommand.Flags().BoolVarP(&deleteAll, "all", "a", false, "delete all tasks")
	rootCmd.AddCommand(deleteCommand)
}

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "for removing the tasks permanently",
	Args: func(cmd *cobra.Command, args []string) error {
		// If deleteAll flag not provided then we need id of tasks
		if deleteAll {
			return nil
		}
		if len(args) < 1 {
			return errors.New("needs id of tasks to be deleted")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Delete all tasks
		if deleteAll {
			err := libs.DeleteAll()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("All tasks deleted")
			return
		}

		err := libs.DeleteByIds(args)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Tasks with ids ", args, " deleted")
	},
}
