package cmd

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/mmilanovic4/orbx/internal/netutil"

	"github.com/spf13/cobra"
)

var headersCmd = &cobra.Command{
	Use:     "headers [url]",
	Short:   "Show HTTP response status and headers",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		// only the headers are needed, the body is never downloaded
		resp, err := netutil.Open(args[0])
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		resp.Body.Close()

		duration := time.Since(start)

		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Time: %s\n", duration)

		// sorted by name, map order would change on every run
		fmt.Println("Headers:")
		for _, key := range slices.Sorted(maps.Keys(resp.Header)) {
			for _, value := range resp.Header[key] {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(headersCmd)
}
