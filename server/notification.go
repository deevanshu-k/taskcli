package server

import (
	"fmt"
	"os/exec"
	"taskcli/structs"
	"time"
)

type NotificationManager struct {
	tasks []structs.Task
}

func NewNotificationManager() *NotificationManager {
	return &NotificationManager{}
}

func (nm *NotificationManager) Start(updatedTasks <-chan []structs.Task) {
	// GET UPDATED TASKS
	go func() {
		for tasks := range updatedTasks {
			nm.tasks = tasks
		}
	}()

	// NOTIFICATION ROUTINE
	go func() {
		for now := range time.Tick(time.Second * 30) {
			for _, task := range nm.tasks {
				fmt.Println(task.NotificationTime, int8(now.Hour()), int8(now.Minute()))
				if task.NotificationTime.IsBefore(int8(now.Hour()), int8(now.Minute())) && task.Notify == structs.ON {
					nm.sendNotification(task)
				}
			}
		}
	}()
}

func (nm *NotificationManager) sendNotification(task structs.Task) {
	// Implementation to send notification for a task
	exec.Command(
		"notify-send",
		"Task Notification",
		fmt.Sprintf("Task: %s\nStatus: %s", task.Name, task.Status.String()),
	).Run()
}
