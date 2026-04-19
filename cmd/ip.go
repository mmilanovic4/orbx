package cmd

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Show public and local IP addresses",
	Run: func(cmd *cobra.Command, args []string) {
		// 🌍 PUBLIC IP
		resp, err := http.Get("https://api.ipify.org")
		if err != nil {
			fmt.Println("Failed to get public IP:", err)
		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			fmt.Println("Public IP:", string(body))
		}

		fmt.Println("\nLocal IPs:")
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

				// filter IPv4 + skip loopback
				if ip == nil || ip.IsLoopback() {
					continue
				}

				ip = ip.To4()
				if ip == nil {
					continue
				}

				fmt.Printf(" - %s (%s)\n", ip.String(), i.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
