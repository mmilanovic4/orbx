package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mmilanovic4/orbx/internal/netutil"

	"github.com/spf13/cobra"
)

var serveHost string

var serveCmd = &cobra.Command{
	Use:     "serve [port]",
	Short:   "Start a static file server in current directory",
	GroupID: "dev",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := netutil.ParsePort(args[0])
		if err != nil {
			return err
		}

		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get directory: %w", err)
		}

		addr := net.JoinHostPort(serveHost, strconv.Itoa(port))

		fmt.Println("Serving:", dir)
		fmt.Printf("On: http://%s\n", addr)

		fs := hideDotfiles(http.FileServer(http.Dir(dir)))

		if err := http.ListenAndServe(addr, fs); err != nil {
			return fmt.Errorf("server error: %w", err)
		}

		return nil
	},
}

// hideDotfiles answers 404 for any path with a segment starting with a dot,
// so files like .env or the .git directory are never served.
func hideDotfiles(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, part := range strings.Split(r.URL.Path, "/") {
			if strings.HasPrefix(part, ".") {
				http.NotFound(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func init() {
	serveCmd.Flags().StringVar(&serveHost, "host", "127.0.0.1", "address to listen on (0.0.0.0 makes it reachable from the network)")
	rootCmd.AddCommand(serveCmd)
}
