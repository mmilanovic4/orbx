package cmd

import (
	"fmt"
	"orbx/internal/netutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:     "ports [port]",
	Short:   "Show processes using network ports",
	GroupID: "dev",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := netutil.ParsePort(args[0])
		if err != nil {
			return err
		}

		out, err := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P").Output()
		if err != nil {
			return fmt.Errorf("failed to run lsof: %w", err)
		}

		lines := strings.Split(string(out), "\n")
		keyword := ":" + strconv.Itoa(port)

		var found []string
		for _, line := range lines {
			if strings.Contains(line, keyword+" ") || strings.Contains(line, keyword+"\n") {
				found = append(found, line)
			}
		}

		if len(found) == 0 {
			fmt.Printf("🟢 Port %d is open.\n", port)
			return nil
		}

		fmt.Printf("🔴 Port %d is in use:\n", port)
		for _, line := range found {
			fmt.Println(line)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
