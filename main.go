package main

import (
	"os"

	"github.com/deevanshu-k/taskcli/cmd"
)

var Version = "v0.0.1"

func main() {
	os.Setenv("APP_VERSION", Version)
	cmd.Execute()
}
