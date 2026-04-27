package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"orbx/internal/sysutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type SSHHost struct {
	Name     string
	HostName string
	User     string
	Port     string
}

var sshlistCmd = &cobra.Command{
	Use:     "sshlist",
	Short:   "List configured SSH hosts from ~/.ssh/config",
	GroupID: "network",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")

		info, err := os.Stat(configPath)
		if err != nil {
			return fmt.Errorf("failed to stat SSH config: %w", err)
		}
		if info.Mode().Perm() != 0600 {
			fmt.Printf("Warning: ~/.ssh/config has permissions %o, should be 600\n\n", info.Mode().Perm())
		}

		data, err := sysutil.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read SSH config: %w", err)
		}

		var hosts []SSHHost
		var current *SSHHost

		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.ToLower(strings.TrimSpace(parts[0]))
			val := strings.TrimSpace(parts[1])

			switch key {
			case "host":
				if current != nil {
					hosts = append(hosts, *current)
				}
				current = &SSHHost{Name: val}
			case "hostname":
				if current != nil {
					current.HostName = val
				}
			case "user":
				if current != nil {
					current.User = val
				}
			case "port":
				if current != nil {
					current.Port = val
				}
			}
		}

		if current != nil {
			hosts = append(hosts, *current)
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read SSH config: %w", err)
		}

		if len(hosts) == 0 {
			fmt.Println("No hosts found in SSH config.")
			return nil
		}

		for i, h := range hosts {
			user := h.User
			if user == "" {
				user = os.Getenv("USER")
			}

			host := h.HostName
			if host == "" {
				host = h.Name
			}

			if h.Port != "" {
				fmt.Printf("%d. %s@%s:%s\n", i+1, user, host, h.Port)
			} else {
				fmt.Printf("%d. %s@%s\n", i+1, user, host)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sshlistCmd)
}
