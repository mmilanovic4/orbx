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

func pathSize(target string) (int64, error) {
	info, err := os.Stat(target)
	if err != nil {
		return 0, fmt.Errorf("failed to access path: %w", err)
	}
	if info.IsDir() {
		return dirSize(target)
	}
	return info.Size(), nil
}

var sizeCmd = &cobra.Command{
	Use:     "size [path...]",
	Short:   "Show logical size of a file or directory",
	GroupID: "util",
	Args:    cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			args = []string{"."}
		}

		if len(args) == 1 {
			size, err := pathSize(args[0])
			if err != nil {
				return err
			}
			abs, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}
			fmt.Printf("%s\t%s\n", formatutil.FormatLogicalSize(size), abs)
			return nil
		}

		var total int64
		for _, target := range args {
			size, err := pathSize(target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: %s\n", err)
				continue
			}
			abs, err := filepath.Abs(target)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}
			fmt.Printf("%s\t%s\n", formatutil.FormatLogicalSize(size), abs)
			total += size
		}

		fmt.Printf("%s\ttotal\n", formatutil.FormatLogicalSize(total))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sizeCmd)
}
