package cmd

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
)

var tcpCmd = &cobra.Command{
	Use:     "tcp [host:port]",
	Short:   "Check TCP port connectivity",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		start := time.Now()

		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		latency := time.Since(start)

		if err != nil {
			fmt.Printf("FAIL %s (%s)\n", target, latency)
			return
		}

		defer conn.Close()

		fmt.Printf("OK   %s (%s)\n", target, latency)
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)
}
