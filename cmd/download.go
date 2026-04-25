package cmd

import (
	"fmt"
	"orbx/internal/netutil"
	"orbx/internal/sysutil"

	"github.com/spf13/cobra"
)

var (
	downloadFile string
)

var downloadCmd = &cobra.Command{
	Use:     "download [url]",
	Short:   "Download a file from a URL",
	GroupID: "util",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := netutil.Get(args[0])
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		if downloadFile != "" {
			if err := sysutil.WriteFile(downloadFile, resp.Body); err != nil {
				return fmt.Errorf("failed to save file: %w", err)
			}
			fmt.Printf("Saved to %s\n", downloadFile)
			return nil
		}

		fmt.Println(string(resp.Body))
		return nil
	},
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadFile, "file", "f", "", "Output file path")
	rootCmd.AddCommand(downloadCmd)
}
