package cmd

import (
	"fmt"
	"orbx/internal/netutil"
	"time"

	"github.com/spf13/cobra"
)

var headersCmd = &cobra.Command{
	Use:     "headers [url]",
	Short:   "Show HTTP response status and headers",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		resp, err := netutil.Get(args[0])
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		duration := time.Since(start)

		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Time: %s\n", duration)

		fmt.Println("Headers:")
		for key, values := range resp.Headers {
			for _, value := range values {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(headersCmd)
}
