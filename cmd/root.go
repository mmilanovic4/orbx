package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version string = "v0.1.0"
var showVersion bool

var rootCmd = &cobra.Command{
	Use:   "orbx",
	Short: "System utility CLI",
	Long:  `orbx is a lightweight CLI utility for quick system tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println("orbx", version)
			return
		}

		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version")
}
