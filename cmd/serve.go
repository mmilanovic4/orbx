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
	Run: func(cmd *cobra.Command, args []string) {
		port, err := netutil.ParsePort(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get directory:", err)
			return
		}

		fmt.Println("Serving:", dir)
		fmt.Printf("On: http://localhost:%d\n", port)

		fs := http.FileServer(http.Dir(dir))

		err = http.ListenAndServe(":"+strconv.Itoa(port), fs)
		if err != nil {
			fmt.Println("Server error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
