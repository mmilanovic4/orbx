package cmd

import (
	"fmt"
	"net"
	"orbx/internal/netutil"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:     "ip",
	Short:   "Show public and local IP addresses",
	GroupID: "network",
	Run: func(cmd *cobra.Command, args []string) {
		// Public IP
		resp, err := netutil.Get("https://api.ipify.org")
		if err != nil {
			fmt.Println("Failed to get public IP:", err)
		} else {
			fmt.Println("Public IP:", string(resp.Body))
		}

		// Local IPs
		fmt.Println("Local IPs:")
		interfaces, err := net.Interfaces()
		if err != nil {
			fmt.Println("Failed to get network interfaces:", err)
			return
		}

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

				// skip loopback
				if ip.IsLoopback() {
					// continue
				}

				// MAC address
				mac := i.HardwareAddr.String()
				if mac == "" {
					mac = "No MAC"
				}

				fmt.Printf(" - %-8s %-39s [%s]\n", "("+i.Name+")", ip.String(), mac)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
