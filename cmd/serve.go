package cmd

import (
	"fmt"
	"net/http"
	"orbx/internal/netutil"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

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

		fmt.Println("Serving:", dir)
		fmt.Printf("On: http://localhost:%d\n", port)

		fs := http.FileServer(http.Dir(dir))

		if err := http.ListenAndServe(":"+strconv.Itoa(port), fs); err != nil {
			return fmt.Errorf("server error: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
