package cmd

import (
	"fmt"
	"time"

	"github.com/mmilanovic4/orbx/internal/netutil"

	"github.com/spf13/cobra"
)

var (
	pingCount int
)

var pingCmd = &cobra.Command{
	Use:     "ping [url]",
	Short:   "HTTP latency check (like ping, but for URLs)",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if pingCount < 1 {
			return fmt.Errorf("count must be at least 1")
		}

		var total, minLatency, maxLatency time.Duration
		succeeded := 0

		for i := 1; i <= pingCount; i++ {
			start := time.Now()

			_, err := netutil.Get(args[0], netutil.WithTimeout(3*time.Second))
			latency := time.Since(start)

			if err != nil {
				fmt.Printf("%d: ERROR (%s)\n", i, err)
				continue
			}

			fmt.Printf("%d: %s\n", i, latency)

			if succeeded == 0 || latency < minLatency {
				minLatency = latency
			}
			if latency > maxLatency {
				maxLatency = latency
			}
			total += latency
			succeeded++
		}

		// failed requests have no latency, so they stay out of the stats
		if succeeded == 0 {
			return fmt.Errorf("all %d requests failed", pingCount)
		}

		fmt.Printf("\navg: %s\n", total/time.Duration(succeeded))
		fmt.Printf("min: %s\n", minLatency)
		fmt.Printf("max: %s\n", maxLatency)
		if failed := pingCount - succeeded; failed > 0 {
			fmt.Printf("failed: %d/%d\n", failed, pingCount)
		}

		return nil
	},
}

func init() {
	pingCmd.Flags().IntVarP(&pingCount, "count", "c", 4, "number of requests")
	rootCmd.AddCommand(pingCmd)
}
