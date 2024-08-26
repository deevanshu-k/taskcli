package cmd

import (
	"fmt"
	"os"

	"github.com/deevanshu-k/taskcli/libs"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var statusPending bool = false
var statusInProgress bool = false
var statusComplete bool = false
var hasDate string = ""

func init() {
	listCommand.PersistentFlags().BoolVarP(&statusPending, "pending", "p", false, "with status pending")
	listCommand.PersistentFlags().BoolVarP(&statusInProgress, "inprogress", "i", false, "with status done")
	listCommand.PersistentFlags().BoolVarP(&statusComplete, "complete", "c", false, "with status complete")
	listCommand.PersistentFlags().StringVarP(&hasDate, "date", "d", "", "with status complete")
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Return all the tasks",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		records, err := libs.AllData()
		if err != nil {
			println(err)
			os.Exit(1)
		}
		// Set headers
		table.SetHeader([]string{"Id", "Task", "Status", "CreatedAt"})
		for _, record := range records {
			needtoappend := false
			// Check for date flag
			if hasDate != "" && hasDate != record[3] {
				continue
			}
			// Checks for status flages
			if !statusPending && !statusComplete && !statusInProgress {
				needtoappend = true
			}
			if statusPending && record[2] == fmt.Sprint(int(libs.Pending)) {
				needtoappend = true
			}
			if statusInProgress && record[2] == fmt.Sprint(int(libs.Inprogress)) {
				needtoappend = true
			}
			if statusComplete && record[2] == fmt.Sprint(int(libs.Complete)) {
				needtoappend = true
			}

			if needtoappend {
				var status string
				if record[2] == fmt.Sprint(int(libs.Inprogress)) {
					status = color.GreenString(libs.Inprogress.String())
				} else if record[2] == fmt.Sprint(int(libs.Complete)) {
					status = color.YellowString(libs.Complete.String())
				} else {
					status = color.RedString(libs.Pending.String())
				}
				table.Append([]string{record[0], record[1], status, record[3]})
			}
		}
		table.Render()
	},
}
