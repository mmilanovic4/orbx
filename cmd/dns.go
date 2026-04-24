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
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordType := ""

		if len(args) > 1 {
			recordType = strings.TrimSpace(strings.ToUpper(args[1]))
		}

		fmt.Println("Resolving:", domain)

		switch recordType {
		case "MX":
			records, err := net.LookupMX(domain)
			if err != nil {
				return fmt.Errorf("failed to lookup MX records: %w", err)
			}
			fmt.Println("\nMX records:")
			for _, mx := range records {
				fmt.Printf("  %s (priority %d)\n", mx.Host, mx.Pref)
			}
		case "CNAME":
			record, err := net.LookupCNAME(domain)
			if err != nil {
				return fmt.Errorf("failed to lookup CNAME record: %w", err)
			}
			fmt.Println("\nCNAME record:")
			fmt.Println(" ", record)
		case "TXT":
			records, err := net.LookupTXT(domain)
			if err != nil {
				return fmt.Errorf("failed to lookup TXT records: %w", err)
			}
			fmt.Println("\nTXT records:")
			for _, t := range records {
				fmt.Println(" ", t)
			}
		default:
			ips, err := net.LookupIP(domain)
			if err != nil {
				return fmt.Errorf("failed to lookup A/AAAA records: %w", err)
			}
			fmt.Println("\nA / AAAA records:")
			for _, ip := range ips {
				fmt.Println(" ", ip)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
