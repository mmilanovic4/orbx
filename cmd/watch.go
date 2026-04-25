package cmd

import (
	"fmt"
	"orbx/internal/sysutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:                "watch [interval] [command...]",
	Short:              "Repeatedly run a command every N seconds",
	GroupID:            "util",
	Args:               cobra.MinimumNArgs(2),
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		interval, err := strconv.ParseFloat(args[0], 64)
		if err != nil || interval <= 0 {
			return fmt.Errorf("invalid interval %q: must be a positive number", args[0])
		}

		command := args[1:]
		ticker := time.NewTicker(time.Duration(interval * float64(time.Second)))
		defer ticker.Stop()

		run := func() {
			sysutil.ClearScreen()
			fmt.Printf("Every %.1fs: %s\n\n", interval, strings.Join(command, " "))

			c := exec.Command(command[0], command[1:]...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Run()
		}

		run()
		for range ticker.C {
			run()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
