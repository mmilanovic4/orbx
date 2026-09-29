package cmd

import (
	"fmt"
	"net"

	"github.com/mmilanovic4/orbx/internal/netutil"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:     "ip",
	Short:   "Show public and local IP addresses",
	GroupID: "network",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := netutil.Get("https://api.ipify.org")
		switch {
		case err != nil:
			fmt.Println("Failed to get public IP:", err)
		case !resp.OK():
			fmt.Println("Failed to get public IP: server responded with", resp.Status)
		default:
			fmt.Println("Public IP:", string(resp.Body))
		}

		fmt.Println("Local IPs:")
		interfaces, err := net.Interfaces()
		if err != nil {
			return fmt.Errorf("failed to get network interfaces: %w", err)
		}

		type row struct {
			name, ip, mac string
		}
		var rows []row

		for _, i := range interfaces {
			addrs, err := i.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				var ip net.IP

				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				if ip.IsLoopback() {
					// continue
				}

				mac := i.HardwareAddr.String()
				if mac == "" {
					mac = "No MAC"
				}

				rows = append(rows, row{"[" + i.Name + "]", ip.String(), mac})
			}
		}

		// interface names like [bridge100] overflowed a fixed column
		nameWidth, ipWidth := 0, 0
		for _, r := range rows {
			nameWidth = max(nameWidth, len(r.name))
			ipWidth = max(ipWidth, len(r.ip))
		}
		for _, r := range rows {
			fmt.Printf(" - %-*s %-*s [%s]\n", nameWidth, r.name, ipWidth, r.ip, r.mac)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
