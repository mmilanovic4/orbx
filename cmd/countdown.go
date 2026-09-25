package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var countdownCmd = &cobra.Command{
	Use:     "countdown [duration]",
	Short:   "Countdown timer (e.g. 1h30m, 5m, 90s)",
	GroupID: "util",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		duration, err := time.ParseDuration(args[0])
		if err != nil {
			return fmt.Errorf("invalid duration %q — examples: 1h30m, 5m, 90s", args[0])
		}

		end := time.Now().Add(duration)
		for {
			remaining := time.Until(end)
			if remaining <= 0 {
				fmt.Print("\rDone.   \n")
				fmt.Print("\a")
				return nil
			}

			// round up, so 3s starts at 00:00:03 and 00:00:00 is never shown
			secs := int((remaining + time.Second - 1) / time.Second)
			h := secs / 3600
			m := secs / 60 % 60
			s := secs % 60

			fmt.Printf("\r%02d:%02d:%02d", h, m, s)

			// sleep until the shown value changes rather than a fixed second,
			// which would slowly drift and skip values
			time.Sleep(remaining - time.Duration(secs-1)*time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(countdownCmd)
}
