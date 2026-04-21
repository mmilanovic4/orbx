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
	Run: func(cmd *cobra.Command, args []string) {
		interval, err := strconv.ParseFloat(args[0], 64)
		if err != nil || interval <= 0 {
			fmt.Println("Invalid interval.")
			return
		}

		command := args[1:]
		ticker := time.NewTicker(time.Duration(interval * float64(time.Second)))
		defer ticker.Stop()

		run := func() {
			sysutil.ClearScreen()
			fmt.Printf("Every %.1fs: %s\n\n", interval, strings.Join(command, " "))

			var c *exec.Cmd
			c = exec.Command(command[0], command[1:]...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Run()
		}

		run()
		for range ticker.C {
			run()
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
