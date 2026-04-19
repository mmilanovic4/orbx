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
	Run: func(cmd *cobra.Command, args []string) {
		ip := strings.TrimSpace(args[0])

		names, err := net.LookupAddr(ip)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(names) == 0 {
			fmt.Println("Hostname not found")
			return
		}

		host := strings.TrimSuffix(names[0], ".")
		fmt.Println(host)
	},
}

func init() {
	rootCmd.AddCommand(rdnsCmd)
}
