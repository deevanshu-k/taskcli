package cmd

import (
	"fmt"
	"strconv"
	"taskcli/structs"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	updateCommand.Flags().StringP("description", "d", "", "Description of the task")
	updateCommand.Flags().StringP("status", "s", "", "Status of the task (P/I/C)")
	updateCommand.Flags().StringP("time", "t", "", "Notification time in HH:MM format")
	updateCommand.Flags().StringP("notify", "n", "", "Task notification (on/off)")

	rootCommand.AddCommand(updateCommand)
}

var updateCommand = &cobra.Command{
	Use:     "update",
	Short:   "Update existing tasks",
	Example: "taskcli update <taskid> \n -d 'New task desc' \n -s 'P/I/C' for updating status P=Pending I=InProgress C=Completed \n -t for setting notification time (HH:MM) \n -n for enabling/disabling notifications (on/off)",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieving flags
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")
		notify, _ := cmd.Flags().GetString("notify")
		var hour int8 = -1
		var minute int8 = -1
		if notificationTime, _ := cmd.Flags().GetString("time"); notificationTime != "" {
			if len(notificationTime) != 5 || notificationTime[2] != ':' {
				fmt.Println("Invalid notification time format")
				return
			}
			t, err := time.Parse("15:04", notificationTime)
			if err != nil {
				fmt.Println("Invalid notification time format")
				return
			}
			hour = int8(t.Hour())
			minute = int8(t.Minute())
		}

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

		if description == "" && (len(status) == 0 || len(status) > 1 || status[0] != 'P' && status[0] != 'I' && status[0] != 'C') && (hour <= -1 || minute <= -1) && (notify != "on" && notify != "off") {
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

		var notificationTime *structs.NotificationTime = nil
		if hour >= 0 && minute >= 0 {
			notificationTime = &structs.NotificationTime{Hour: hour, Minute: minute}
			notify = "on"
		}

		var pnotify *structs.Notify = nil
		if notify == "on" {
			temp := structs.ON
			pnotify = &temp
		}
		if notify == "off" {
			temp := structs.OFF
			pnotify = &temp
		}

		command := structs.NewCommand(structs.UPDATE, &taskId, pdescription, rstatus, nil, notificationTime, pnotify)

		res, err := command.SendCommand()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Println(res)
	},
}
