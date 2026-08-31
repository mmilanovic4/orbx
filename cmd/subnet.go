package cmd

import (
	"fmt"
	"github.com/mmilanovic4/orbx/internal/netutil"

	"github.com/spf13/cobra"
)

var subnetCmd = &cobra.Command{
	Use:   "subnet <ip/cidr> [netmask]",
	Short: "Calculate subnet details (network, broadcast, host range)",
	Long: `Calculate subnet details from an IP and CIDR or netmask.

Usage:
  orbx subnet 192.168.1.10/24
  orbx subnet 192.168.1.10 255.255.255.0
  orbx subnet 192.168.1.10 0xffffff00
  orbx subnet 2001:db8::/32`,
	GroupID: "network",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mask := ""
		if len(args) == 2 {
			mask = args[1]
		}

		info, err := netutil.ParseSubnet(args[0], mask)
		if err != nil {
			return err
		}

		fmt.Printf("%-14s%s\n", "Address:", info.Address)
		fmt.Printf("%-14s%s\n", "Network:", info.Network)
		fmt.Printf("%-14s%s\n", "Netmask:", info.Netmask)
		if !info.IsIPv6 {
			fmt.Printf("%-14s%s\n", "Wildcard:", info.Wildcard)
			fmt.Printf("%-14s%s\n", "Broadcast:", info.Broadcast)
		}
		fmt.Printf("%-14s%s - %s\n", "Host range:", info.FirstHost, info.LastHost)
		fmt.Printf("%-14s%s\n", "Usable hosts:", info.UsableHosts)
		fmt.Printf("%-14s%s\n", "Total hosts:", info.TotalHosts)
		if !info.IsIPv6 {
			fmt.Printf("%-14s%s\n", "Class:", info.Class)
		}
		fmt.Printf("%-14s%t\n", "Private:", info.Private)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(subnetCmd)
}
