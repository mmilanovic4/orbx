package cmd

import (
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is set for release builds (see .goreleaser.yaml) with
// -ldflags "-X github.com/mmilanovic4/orbx/cmd.version=v1.2.0"
var version string

// buildVersion falls back to the version Go records in the binary: the tag
// for "go install github.com/mmilanovic4/orbx@v1.2.0" and for builds of a
// tagged checkout, a pseudo-version for other commits.
func buildVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

var rootCmd = &cobra.Command{
	Use:     "orbx",
	Short:   "System utility CLI",
	Long:    "orbx is a lightweight CLI utility for quick system tasks.",
	Version: buildVersion(),
	// Runs after flag and argument validation, so usage is still shown for
	// usage mistakes but not for errors that happen while a command runs.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
	},
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
}
