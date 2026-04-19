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
	GroupID: "system",
	Run: func(cmd *cobra.Command, args []string) {
		switch runtime.GOOS {
		case "darwin":
			// macOS
			cmd := exec.Command("pbcopy")
			cmd.Stdin = nil
			if err := cmd.Run(); err != nil {
				fmt.Println("error clearing clipboard:", err)
				return
			}
		case "linux":
			// Linux (X11)
			cmd := exec.Command("bash", "-c", "echo -n | xclip -selection clipboard")
			if err := cmd.Run(); err != nil {
				fmt.Println("error clearing clipboard (xclip missing?):", err)
				return
			}
		default:
			fmt.Println("unsupported OS:", runtime.GOOS)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(clearclipCmd)
}
