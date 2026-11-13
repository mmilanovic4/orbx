package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mmilanovic4/orbx/internal/netutil"
	"github.com/mmilanovic4/orbx/internal/sysutil"

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
		resp, err := netutil.Open(args[0])
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("download failed: server responded with %s", resp.Status)
		}

		// the body is streamed, so files larger than the available memory
		// download as well, and an interrupted download never replaces an
		// existing file with a partial one
		if downloadFile != "" {
			if err := sysutil.ReplaceFileFrom(downloadFile, resp.Body); err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
			fmt.Printf("Saved to %s\n", downloadFile)
			return nil
		}

		if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		return nil
	},
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadFile, "file", "f", "", "Output file path")
	rootCmd.AddCommand(downloadCmd)
}
