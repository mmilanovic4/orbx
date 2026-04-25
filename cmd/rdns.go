package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

var rdnsCmd = &cobra.Command{
	Use:     "rdns [ip]",
	Short:   "Reverse DNS lookup for an IP address",
	GroupID: "network",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip := strings.TrimSpace(args[0])

		if net.ParseIP(ip) == nil {
			return fmt.Errorf("invalid IP address: %s", ip)
		}

		names, err := net.LookupAddr(ip)
		if err != nil {
			return fmt.Errorf("lookup failed: %w", err)
		}

		if len(names) == 0 {
			return fmt.Errorf("no hostname found for %s", ip)
		}

		fmt.Println(strings.TrimSuffix(names[0], "."))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rdnsCmd)
}
