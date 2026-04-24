package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var clearclipCmd = &cobra.Command{
	Use:     "clearclip",
	Short:   "Clear system clipboard",
	GroupID: "util",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch runtime.GOOS {
		case "darwin":
			// macOS
			c := exec.Command("pbcopy")
			c.Stdin = nil
			if err := c.Run(); err != nil {
				return fmt.Errorf("error clearing clipboard: %w", err)
			}
		case "linux":
			// Linux (X11)
			c := exec.Command("xclip", "-selection", "clipboard")
			c.Stdin = nil
			if err := c.Run(); err != nil {
				return fmt.Errorf("error clearing clipboard (xclip missing?): %w", err)
			}
		default:
			return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearclipCmd)
}
