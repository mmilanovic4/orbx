package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	VERSION string = "v0.1.0"
)

var rootCmd = &cobra.Command{
	Use:     "orbx",
	Short:   "System utility CLI",
	Long:    "orbx is a lightweight CLI utility for quick system tasks.",
	Version: VERSION,
	Run: func(cmd *cobra.Command, args []string) {
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
