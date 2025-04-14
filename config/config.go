package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var Config *TaskCliConfig

type TaskCliConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Load() {
	viper.SetConfigName("taskcli")
	viper.SetConfigType("json")
	viper.AddConfigPath(ConfigDir)

	configFilePath := filepath.Join(ConfigDir, "taskcli.json")
	storageFilePath := filepath.Join(StorageDir, "taskcli.db")

	// Check if the config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		createDefaultConfig(configFilePath)
	}

	// Check if the storage file exists
	if _, err := os.Stat(storageFilePath); os.IsNotExist(err) {
		createDefaultStorage(storageFilePath)
	}

	// Now read the config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fatal error reading config file: %v", err)
	}

	// Unmarshal into struct
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}
}

func createDefaultConfig(configFilePath string) {
	// Create parent directory if not exists
	dir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	// Create the default config file
	defaultConfig := TaskCliConfig{
		Host: "127.0.0.1",
		Port: 0,
	}

	viper.SetDefault("host", defaultConfig.Host)
	viper.SetDefault("port", defaultConfig.Port)

	if err := viper.WriteConfigAs(configFilePath); err != nil {
		log.Fatalf("Failed to create default config file: %v", err)
	}
}

func createDefaultStorage(storageFilePath string) {
	// Create the storage file if it doesn't exist
	if err := os.MkdirAll(StorageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}
	if _, err := os.Create(storageFilePath); err != nil {
		log.Fatalf("Failed to create storage file: %v", err)
	}
}
