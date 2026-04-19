package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var pingCount int

var pingCmd = &cobra.Command{
	Use:     "ping [url]",
	Short:   "HTTP latency check (like ping, but for URLs)",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		// ensure scheme exists
		if !strings.HasPrefix(url, "http") {
			url = "https://" + url
		}

		client := &http.Client{
			Timeout: 3 * time.Second,
		}

		var total time.Duration
		var min time.Duration = time.Hour
		var max time.Duration

		fmt.Printf("PING %s\n\n", url)

		for i := 1; i <= pingCount; i++ {

			start := time.Now()

			resp, err := client.Get(url)
			latency := time.Since(start)

			if err != nil {
				fmt.Printf("%d: ERROR (%s)\n", i, err)
				continue
			}

			resp.Body.Close()

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
	},
}

func init() {
	pingCmd.Flags().IntVarP(&pingCount, "count", "c", 4, "number of requests")
	rootCmd.AddCommand(pingCmd)
}
