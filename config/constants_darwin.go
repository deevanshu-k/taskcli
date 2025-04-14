//go:build darwin

package config

import (
	"log"
	"os"
	"path/filepath"
)

var (
	ConfigDir  string
	StorageDir string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Cannot get home directory: %v", err)
	}

	ConfigDir = filepath.Join(home, "Library", "Application Support", "taskcli")
	StorageDir = filepath.Join(home, "Library", "Application Support", "taskcli")
}
