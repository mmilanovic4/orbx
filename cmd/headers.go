package cmd

import (
	"fmt"
	"net/http"
	"orbx/internal/netutil"
	"time"

	"github.com/spf13/cobra"
)

var headersCmd = &cobra.Command{
	Use:     "headers [url]",
	Short:   "Show HTTP response status and headers",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := netutil.NormalizeURL(args[0])

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		start := time.Now()

		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Request failed:", err)
			return
		}
		defer resp.Body.Close()

		duration := time.Since(start)

		// Status
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Time: %s\n\n", duration)

		// Headers
		fmt.Println("Headers:")
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(headersCmd)
}
