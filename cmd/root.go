package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var showVersion bool

var rootCmd = &cobra.Command{
	Use:   "taskcli",
	Short: "A bare minimum task management through cli.",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println("Version: ", os.Getenv("APP_VERSION"))
		}
	},
}

func Execute() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", true, "Prints version")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
