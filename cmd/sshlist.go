package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mmilanovic4/orbx/internal/netutil"
	"github.com/mmilanovic4/orbx/internal/sysutil"

	"github.com/spf13/cobra"
)

// formatSSHHost shows the alias passed to ssh and where it connects,
// e.g. "myserver → deploy@10.0.0.5:2222 (~/.ssh/id_ed25519)".
func formatSSHHost(h netutil.SSHHost) string {
	s := h.Alias

	if h.HostName != "" || h.User != "" || h.Port != "" {
		target := h.HostName
		if target == "" {
			target = h.Alias
		}
		if h.User != "" {
			target = h.User + "@" + target
		}
		if h.Port != "" {
			target += ":" + h.Port
		}
		s += " → " + target
	}

	if h.IdentityFile != "" {
		s += " (" + h.IdentityFile + ")"
	}

	return s
}

var sshlistCmd = &cobra.Command{
	Use:     "sshlist",
	Short:   "List configured SSH hosts from ~/.ssh/config",
	GroupID: "network",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to find home directory: %w", err)
		}
		configPath := filepath.Join(home, ".ssh", "config")

		info, err := os.Stat(configPath)
		if err != nil {
			return fmt.Errorf("failed to stat SSH config: %w", err)
		}
		// the same check ssh does, it refuses a config others can write to
		if perm := info.Mode().Perm(); perm&0o022 != 0 {
			fmt.Printf("⚠ Warning: ~/.ssh/config has permissions %o, ssh refuses it while group or others can write to it (chmod 600)\n\n", perm)
		}

		data, err := sysutil.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read SSH config: %w", err)
		}

		hosts, err := netutil.ParseSSHConfig(data)
		if err != nil {
			return fmt.Errorf("failed to read SSH config: %w", err)
		}

		if len(hosts) == 0 {
			fmt.Println("No hosts found in SSH config.")
			return nil
		}

		for i, h := range hosts {
			fmt.Printf("%d. %s\n", i+1, formatSSHHost(h))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sshlistCmd)
}
