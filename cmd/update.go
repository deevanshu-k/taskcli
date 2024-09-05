package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/deevanshu-k/taskcli/libs"
	"github.com/spf13/cobra"
)

var updateStatus string = ""
var updateTask string = ""

func init() {
	updateCommand.Flags().StringVarP(&updateStatus, "status", "s", "", "For updating the task status, Values are \n '0' for pending,'1' for inprogress,'2' for complete,")
	updateCommand.Flags().StringVarP(&updateTask, "task", "t", "", "For changing the task")
	rootCmd.AddCommand(updateCommand)
}

var updateCommand = &cobra.Command{
	Use:   "update taskId [-s status|-t task]",
	Short: "Update the task detail",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("only taskId argument is allowed")
		}
		if len(args) == 0 {
			return errors.New("id of task is requeird")
		}
		if updateStatus != "" && updateStatus != "0" && updateStatus != "1" && updateStatus != "2" {
			return errors.New("-s wrong value passed")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if updateStatus != "" {
			id, _ := strconv.Atoi(args[0])
			status, _ := strconv.Atoi(updateStatus)
			err := libs.UpdateStatus(id, status)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if updateTask != "" {
			err := libs.UpdateTask(args[0], updateTask)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		listCommand.Run(listCommand, []string{})
	},
	DisableFlagsInUseLine: true,
}
