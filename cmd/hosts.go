package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"orbx/internal/sysutil"
	"strings"

	"github.com/spf13/cobra"
)

type HostEntry struct {
	IP      string
	Domains []string
}

var hostsCmd = &cobra.Command{
	Use:     "hosts",
	Short:   "List entries from /etc/hosts",
	GroupID: "network",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := sysutil.ReadFile("/etc/hosts")
		if err != nil {
			return fmt.Errorf("failed to read /etc/hosts: %w", err)
		}

		var entries []HostEntry
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}

			entries = append(entries, HostEntry{
				IP:      fields[0],
				Domains: fields[1:],
			})
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read /etc/hosts: %w", err)
		}

		if len(entries) == 0 {
			fmt.Println("No entries found in /etc/hosts.")
			return nil
		}

		for i, e := range entries {
			fmt.Printf("%d. %s → %s\n", i+1, e.IP, strings.Join(e.Domains, ", "))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
