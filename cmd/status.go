package cmd

import (
	"fmt"
	"os"

	"github.com/deevanshu-k/taskcli/libs"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCommand)
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "Give info related to no of task of different catagory",
	Run: func(cmd *cobra.Command, args []string) {
		// Find Stats
		records, err := libs.AllData()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		noOfPendingTask := 0
		noOfInProgressTask := 0
		noOfCompleteTask := 0
		for _, record := range records {
			if record[2] == fmt.Sprint(int(libs.Pending)) {
				noOfPendingTask++
			}
			if record[2] == fmt.Sprint(int(libs.Inprogress)) {
				noOfInProgressTask++
			}
			if record[2] == fmt.Sprint(int(libs.Complete)) {
				noOfCompleteTask++
			}
		}

		// Render Stats
		fmt.Println("Pending: ", color.RedString(fmt.Sprint(noOfPendingTask)),
			"   ", "In-Progress: ", color.GreenString(fmt.Sprint(noOfInProgressTask)),
			"   ", "Complete: ", color.YellowString(fmt.Sprint(noOfCompleteTask)))
	},
	DisableFlagsInUseLine: true,
}
