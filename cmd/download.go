package cmd

import (
	"fmt"
	"orbx/internal/netutil"
	"orbx/internal/sysutil"

	"github.com/spf13/cobra"
)

var downloadFile string

var downloadCmd = &cobra.Command{
	Use:     "download [url]",
	Short:   "Download a file from a URL",
	GroupID: "util",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := netutil.Get(args[0])
		if err != nil {
			fmt.Println("Download failed:", err)
			return
		}

		fmt.Println(string(resp.Body))

		if downloadFile != "" {
			err := sysutil.WriteFile(downloadFile, resp.Body)
			if err != nil {
				fmt.Println("Failed to save file:", err)
				return
			}

			fmt.Printf("Saved to %s\n", downloadFile)
			return
		}
	},
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadFile, "file", "f", "", "Output file path")
	rootCmd.AddCommand(downloadCmd)
}
