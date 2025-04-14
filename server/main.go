package server

import (
	"log/slog"
	"taskcli/config"
)

func printConfiguration() {
	slog.Info("starting TaskCli", slog.String("Version", "1.0.1"))
	slog.Info("running with", slog.String("Host", config.Config.Host))
	slog.Info("running with", slog.Int("Port", config.Config.Port))
}

func Start() {
	printConfiguration()

	server := NewServer(config.Config.Host, config.Config.Port)
	slog.Info("server started", slog.String("Host", config.Config.Host), slog.Int("Port", config.Config.Port))
	server.Start()
}
