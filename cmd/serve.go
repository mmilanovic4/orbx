package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve [port]",
	Short:   "Start a static file server in current directory",
	GroupID: "dev",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		portStr := args[0]

		// strict parse (no stripping)
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println("Invalid port:", portStr)
			return
		}

		// valid port range
		if port < 1 || port > 65535 {
			fmt.Println("Port must be between 1 and 65535")
			return
		}

		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get directory:", err)
			return
		}

		fmt.Println("Serving:", dir)
		fmt.Println("On: http://localhost:" + portStr)

		fs := http.FileServer(http.Dir(dir))

		err = http.ListenAndServe(":"+portStr, fs)
		if err != nil {
			fmt.Println("Server error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
