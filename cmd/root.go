package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	VERSION string = "v0.1.0"
)

var (
	showVersion bool
)

var rootCmd = &cobra.Command{
	Use:   "orbx",
	Short: "System utility CLI",
	Long:  "orbx is a lightweight CLI utility for quick system tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println("orbx", VERSION)
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
	// Version flag
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version")

	// 🧰 Utilities
	rootCmd.AddGroup(&cobra.Group{
		ID:    "util",
		Title: "🧰 Utilities",
	})

	// 🌐 Network Tools
	rootCmd.AddGroup(&cobra.Group{
		ID:    "network",
		Title: "🌐 Network Tools",
	})

	// 💻 Developer Tools
	rootCmd.AddGroup(&cobra.Group{
		ID:    "dev",
		Title: "💻 Developer Tools",
	})

	// 📦 Content
	rootCmd.AddGroup(&cobra.Group{
		ID:    "misc",
		Title: "📦 Content",
	})
}
