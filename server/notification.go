package server

import (
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
	"taskcli/config"
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
		for now := range time.Tick(time.Second * time.Duration(config.Config.NotificationFrequency)) {
			for _, task := range nm.tasks {
				if task.NotificationTime.IsBefore(int8(now.Hour()), int8(now.Minute())) && task.Notify == structs.ON {
					if err := nm.sendNotification(task); err != nil {
						slog.Error("Error sending notification:", "Error", err, "Required", "linux: notify-send, windows: New-BurntToastNotification, Darwin: osascript")
					}
				}
			}
		}
	}()
}

func (nm *NotificationManager) sendNotification(task structs.Task) error {
	title := "🔔 Task Reminder"
	message := task.Name

	switch runtime.GOOS {
	case "linux":
		// Linux: notify-send
		return exec.Command("notify-send", title, message).Run()

	case "darwin":
		// macOS: osascript (AppleScript)
		cmd := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
		return exec.Command("osascript", "-e", cmd).Run()

	case "windows":
		// Windows: PowerShell toast notification
		powershell := fmt.Sprintf(`New-BurntToastNotification -Text '%s','%s'`, title, message)
		return exec.Command("powershell", "-Command", powershell).Run()

	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
