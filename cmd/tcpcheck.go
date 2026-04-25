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

		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		latency := time.Since(start)

		if err != nil {
			fmt.Printf("🔴 %s (%s)\n", target, latency)
			return nil
		}
		defer conn.Close()

		fmt.Printf("🟢 %s (%s)\n", target, latency)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tcpcheckCmd)
}
