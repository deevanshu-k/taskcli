package server

import (
	"log/slog"
	"taskcli/config"
)

func printConfiguration() {
	slog.Info("starting TaskCli", slog.String("Version", "1.0.1"))
	slog.Info("running with", slog.String("Host", config.Config.Host))
	slog.Info("running with", slog.Int("Port", config.Config.Port))
	slog.Info("running with", slog.Int("NotificationFrequency", config.Config.NotificationFrequency))
}

func Start() {
	printConfiguration()

	// CREATE TCP SERVER
	server := NewServer(config.Config.Host, config.Config.Port)

	// START NOTIFICATION ROUTINE
	notificationManager := NewNotificationManager()
	notificationManager.Start(server.Manager.GetUpdatedTasks())
	slog.Info("notification manager started", slog.Int("Frequency", config.Config.NotificationFrequency))

	// START TCP SERVER
	slog.Info("server started", slog.String("Host", config.Config.Host), slog.Int("Port", config.Config.Port))
	server.Start()
}
