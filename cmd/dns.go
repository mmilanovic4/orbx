package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:     "dns [domain] [type]",
	Short:   "Resolve DNS records for a domain",
	GroupID: "network",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		recordType := ""

		if len(args) > 1 {
			recordType = strings.TrimSpace(strings.ToUpper(args[1]))
		}

		fmt.Println("Resolving:", domain)

		switch recordType {
		case "MX":
			// MX records
			mxRecords, err := net.LookupMX(domain)
			if err == nil && len(mxRecords) > 0 {
				fmt.Println("\nMX records:")
				for _, mx := range mxRecords {
					fmt.Printf("  %s (priority %d)\n", mx.Host, mx.Pref)
				}
			}
		case "CNAME":
			// CNAME record
			cnameRecord, err := net.LookupCNAME(domain)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("\nCNAME record:")
			fmt.Println(" ", cnameRecord)
		case "TXT":
			// TXT records
			txtRecords, err := net.LookupTXT(domain)
			if err != nil && len(txtRecords) > 0 {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("\nTXT records:")
			for _, t := range txtRecords {
				fmt.Println(" ", t)
			}
		default:
			// A / AAAA records
			ips, err := net.LookupIP(domain)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("\nA / AAAA records:")
			for _, ip := range ips {
				fmt.Println(" ", ip)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
