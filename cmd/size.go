package cmd

import (
	"fmt"
	"orbx/internal/formatutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})
	return size, err
}

var sizeCmd = &cobra.Command{
	Use:     "size [path]",
	Short:   "Show logical size of a file or directory",
	GroupID: "util",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "."
		if len(args) > 0 {
			target = args[0]
		}

		info, err := os.Lstat(target)
		if err != nil {
			return fmt.Errorf("failed to access path: %w", err)
		}

		abs, err := filepath.Abs(target)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}

		var totalSize int64
		if info.IsDir() {
			totalSize, err = dirSize(target)
			if err != nil {
				return fmt.Errorf("failed to calculate size: %w", err)
			}
		} else {
			totalSize = info.Size()
		}

		fmt.Printf("%s\t%s\n", formatutil.FormatLogicalSize(totalSize), abs)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sizeCmd)
}
