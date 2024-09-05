package main

import (
	"os"

	"github.com/deevanshu-k/taskcli/cmd"
)

var Version = "v0.0.7"
var BaseUrl = "http://localhost:5000"

func main() {
	os.Setenv("APP_VERSION", Version)
	os.Setenv("BASE_URL", BaseUrl)
	cmd.Execute()
}
