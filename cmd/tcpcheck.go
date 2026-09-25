package cmd

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
)

var tcpcheckCmd = &cobra.Command{
	Use:     "tcpcheck [host:port]",
	Short:   "Check TCP port connectivity",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		if _, port, err := net.SplitHostPort(target); err != nil || port == "" {
			return fmt.Errorf("invalid target %q: expected host:port, e.g. example.com:443 or [::1]:22", target)
		}

		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		latency := time.Since(start)

		if err != nil {
			fmt.Printf("🔴 %s (%s)\n", target, latency)
			// a non-zero exit code lets scripts use it: orbx tcpcheck db:5432 && ...
			return fmt.Errorf("connection failed: %w", err)
		}
		defer conn.Close()

		fmt.Printf("🟢 %s (%s)\n", target, latency)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tcpcheckCmd)
}
