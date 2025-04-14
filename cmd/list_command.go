package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"taskcli/structs"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	listCommand.Flags().StringP("filter", "f", "", "Filter tasks by status")

	rootCommand.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:     "list",
	Short:   "List all tasks",
	Example: "taskcli list",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")
		var r *structs.Status = nil

		if len(filter) > 0 && rune(filter[0]) == rune(structs.PENDING) {
			temp := structs.PENDING
			r = &temp
		}
		if len(filter) > 0 && rune(filter[0]) == rune(structs.INPROGRESS) {
			temp := structs.INPROGRESS
			r = &temp
		}
		if len(filter) > 0 && rune(filter[0]) == rune(structs.COMPLETED) {
			temp := structs.COMPLETED
			r = &temp
		}

		command := structs.NewCommand(structs.LIST, nil, nil, r, nil)

		res, err := command.SendCommand()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		sort.Slice(res.Tasks, func(i, j int) bool {
			return res.Tasks[i].Id < res.Tasks[j].Id
		})

		printTasks(res.Tasks)
	},
}

func printTasks(tasks []structs.Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Status"})

	for _, t := range tasks {
		table.Append([]string{strconv.Itoa(t.Id), t.Name, t.Status.String(), t.Date})
	}

	table.Render() // Print it!
}
