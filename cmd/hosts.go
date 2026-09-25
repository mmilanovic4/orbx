package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/mmilanovic4/orbx/internal/sysutil"

	"github.com/spf13/cobra"
)

type HostEntry struct {
	IP      string
	Domains []string
}

func parseHosts(data []byte) ([]HostEntry, error) {
	var entries []HostEntry
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()

		// everything from # to the end of the line is a comment,
		// including after an entry
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = line[:i]
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
		return nil, err
	}
	return entries, nil
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

		entries, err := parseHosts(data)
		if err != nil {
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
