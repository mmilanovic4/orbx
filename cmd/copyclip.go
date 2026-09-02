package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/mmilanovic4/orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var copyclipFile string

var copyclipCmd = &cobra.Command{
	Use:     "copyclip [input]",
	Short:   "Copy input to system clipboard",
	GroupID: "util",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, copyclipFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		switch runtime.GOOS {
		case "darwin":
			c := exec.Command("pbcopy")
			c.Stdin = bytes.NewReader(data)
			if err := c.Run(); err != nil {
				return fmt.Errorf("error copying to clipboard: %w", err)
			}
		case "linux":
			c := exec.Command("xclip", "-selection", "clipboard")
			c.Stdin = bytes.NewReader(data)
			if err := c.Run(); err != nil {
				return fmt.Errorf("error copying to clipboard (xclip missing?): %w", err)
			}
		default:
			return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
		}

		return nil
	},
}

func init() {
	copyclipCmd.Flags().StringVarP(&copyclipFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(copyclipCmd)
}
