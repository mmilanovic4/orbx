package cmd

import (
	"context"
	"fmt"
	"net"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var dnsRecordTypes = []string{"A", "AAAA", "MX", "CNAME", "TXT"}

var dnsCmd = &cobra.Command{
	Use:   "dns [domain] [type]",
	Short: "Resolve DNS records for a domain",
	Long: `Resolve DNS records for a domain.

Supported record types:
  A, AAAA  (both by default)
  MX
  CNAME
  TXT`,
	GroupID: "network",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordType := ""

		if len(args) > 1 {
			recordType = strings.TrimSpace(strings.ToUpper(args[1]))
			if !slices.Contains(dnsRecordTypes, recordType) {
				return fmt.Errorf("unsupported record type %q: use %s", args[1], strings.Join(dnsRecordTypes, ", "))
			}
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
			// "ip" asks for both address families, "ip4" and "ip6" for one
			network, label := "ip", "A / AAAA"
			switch recordType {
			case "A":
				network, label = "ip4", "A"
			case "AAAA":
				network, label = "ip6", "AAAA"
			}

			ips, err := net.DefaultResolver.LookupIP(context.Background(), network, domain)
			if err != nil {
				return fmt.Errorf("failed to lookup %s records: %w", label, err)
			}
			fmt.Printf("\n%s records:\n", label)
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
