package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mmilanovic4/orbx/internal/sysutil"

	"github.com/spf13/cobra"
)

// same lower bound as watch(1), anything shorter just spins the CPU
const minWatchInterval = 0.1

var watchCmd = &cobra.Command{
	Use:   "watch [interval] [command...]",
	Short: "Repeatedly run a command every N seconds",
	Long: `Repeatedly run a command every N seconds.

Like watch(1), the command is run with "sh -c" ("cmd /C" on Windows),
so pipes and other shell syntax work when the command is quoted:

  orbx watch 2 'ls | wc -l'`,
	GroupID: "util",
	Args: func(cmd *cobra.Command, args []string) error {
		if isHelpArg(args) {
			return nil
		}
		return cobra.MinimumNArgs(2)(cmd, args)
	},
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// flag parsing is disabled so the watched command keeps its flags,
		// which leaves --help to be handled here
		if isHelpArg(args) {
			return cmd.Help()
		}

		interval, err := strconv.ParseFloat(args[0], 64)
		if err != nil || math.IsNaN(interval) || math.IsInf(interval, 0) || interval < minWatchInterval {
			return fmt.Errorf("invalid interval %q: must be a number of seconds, at least %g", args[0], minWatchInterval)
		}

		command := strings.Join(args[1:], " ")
		ticker := time.NewTicker(time.Duration(interval * float64(time.Second)))
		defer ticker.Stop()

		run := func() error {
			sysutil.ClearScreen()
			fmt.Printf("Every %.1fs: %s\n\n", interval, command)

			c := shellCommand(command)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			// a non-zero exit is part of what is being watched and the shell
			// has already printed why, any other error means nothing ran
			var exitErr *exec.ExitError
			if err := c.Run(); err != nil && !errors.As(err, &exitErr) {
				return fmt.Errorf("failed to run command: %w", err)
			}
			return nil
		}

		if err := run(); err != nil {
			return err
		}
		for range ticker.C {
			if err := run(); err != nil {
				return err
			}
		}

		return nil
	},
}

func isHelpArg(args []string) bool {
	return len(args) > 0 && (args[0] == "-h" || args[0] == "--help")
}

func shellCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", command)
	}
	return exec.Command("sh", "-c", command)
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
