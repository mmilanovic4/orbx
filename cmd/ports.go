package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/mmilanovic4/orbx/internal/netutil"

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

		listeners, err := portListeners(port)
		if err != nil {
			return err
		}

		if len(listeners) > 0 {
			fmt.Printf("🔴 Port %d is in use:\n", port)
			for _, line := range listeners {
				fmt.Println(line)
			}
			return nil
		}

		// lsof only lists the processes it may see, so the port can still
		// be taken by another user's process (or lsof may be missing)
		if portTaken(port) {
			fmt.Printf("🔴 Port %d is in use, by a process that is not visible to this user", port)
			if runtime.GOOS != "windows" && os.Geteuid() != 0 {
				fmt.Print(" (run with sudo to see it)")
			}
			fmt.Println()
			return nil
		}

		fmt.Printf("🟢 Port %d is free.\n", port)
		return nil
	},
}

// portListeners returns the lsof lines (header first) for processes
// listening on TCP port, or nothing if there are none or lsof is missing.
func portListeners(port int) ([]string, error) {
	if _, err := exec.LookPath("lsof"); err != nil {
		return nil, nil
	}

	out, err := exec.Command("lsof", "-nP", "-iTCP:"+strconv.Itoa(port), "-sTCP:LISTEN").Output()
	if err != nil {
		// lsof exits with 1 when nothing matches
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 && len(out) == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to run lsof: %w", err)
	}

	var lines []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// portTaken tries to listen on port. A wildcard listener alone is not
// enough: on macOS it succeeds even while 127.0.0.1 is taken.
func portTaken(port int) bool {
	p := strconv.Itoa(port)
	for _, addr := range []string{":" + p, "127.0.0.1:" + p, "[::1]:" + p} {
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close()
			continue
		}
		if addrInUse(err) {
			return true
		}
		// anything else (no IPv6, a privileged port) says nothing about
		// whether the port is free
	}
	return false
}

func addrInUse(err error) bool {
	if errors.Is(err, syscall.EADDRINUSE) {
		return true
	}
	var errno syscall.Errno
	return runtime.GOOS == "windows" && errors.As(err, &errno) && errno == 10048 // WSAEADDRINUSE
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
