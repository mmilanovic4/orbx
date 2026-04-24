package cmd

import (
	"fmt"
	"orbx/internal/netutil"
	"time"

	"github.com/spf13/cobra"
)

var pingCount int

var pingCmd = &cobra.Command{
	Use:     "ping [url]",
	Short:   "HTTP latency check (like ping, but for URLs)",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var total time.Duration
		var min time.Duration = time.Hour
		var max time.Duration

		for i := 1; i <= pingCount; i++ {
			start := time.Now()

			_, err := netutil.Get(args[0], netutil.WithTimeout(3*time.Second))
			latency := time.Since(start)

			if err != nil {
				fmt.Printf("%d: ERROR (%s)\n", i, err)
				continue
			}

			fmt.Printf("%d: %s\n", i, latency)

			total += latency

			if latency < min {
				min = latency
			}
			if latency > max {
				max = latency
			}
		}

		avg := total / time.Duration(pingCount)

		fmt.Printf("\navg: %s\n", avg)
		fmt.Printf("min: %s\n", min)
		fmt.Printf("max: %s\n", max)

		return nil
	},
}

func init() {
	pingCmd.Flags().IntVarP(&pingCount, "count", "c", 4, "number of requests")
	rootCmd.AddCommand(pingCmd)
}
