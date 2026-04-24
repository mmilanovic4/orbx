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

			h := int(remaining.Hours())
			m := int(remaining.Minutes()) % 60
			s := int(remaining.Seconds()) % 60

			fmt.Printf("\r%02d:%02d:%02d", h, m, s)
			time.Sleep(time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(countdownCmd)
}
