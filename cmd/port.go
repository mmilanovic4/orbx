package cmd

import (
	"fmt"
	"orbx/internal/netutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:     "port [number]",
	Short:   "Show processes using network ports",
	GroupID: "dev",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port, err := netutil.ParsePort(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		out, err := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P").Output()
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		lines := strings.Split(string(out), "\n")

		var found []string
		var keyword string = ":" + strconv.Itoa(port)

		for _, line := range lines {
			if strings.Contains(line, keyword+" ") || strings.Contains(line, keyword+"\n") {
				found = append(found, line)
			}
		}

		if len(found) == 0 {
			fmt.Printf("🟢 Port %d is open.\n", port)
			return
		}

		fmt.Printf("🔴 Port %d is in use:\n", port)

		for _, line := range found {
			fmt.Println(line)
		}
	},
}

func init() {
	rootCmd.AddCommand(portCmd)
}
