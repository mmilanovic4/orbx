package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns [domain]",
	Short: "Resolve DNS records for a domain",
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]

		fmt.Println("Resolving:", domain)

		// 🌐 A records
		ips, err := net.LookupIP(domain)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("\nA / AAAA records:")
		for _, ip := range ips {
			fmt.Println(" ", ip)
		}

		// 📧 MX records
		mxRecords, err := net.LookupMX(domain)
		if err == nil && len(mxRecords) > 0 {
			fmt.Println("\nMX records:")
			for _, mx := range mxRecords {
				fmt.Printf("  %s (priority %d)\n", mx.Host, mx.Pref)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
