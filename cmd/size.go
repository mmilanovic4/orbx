package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mmilanovic4/orbx/internal/formatutil"

	"github.com/spf13/cobra"
)

func dirSize(root string) (int64, error) {
	var size int64
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if path == root {
				return err
			}
			// one unreadable entry should not hide the size of everything else
			fmt.Fprintf(os.Stderr, "warning: %s\n", err)
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: %s\n", err)
			return nil
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

		type row struct {
			size  string
			label string
		}
		var rows []row

		var total int64
		for _, target := range args {
			size, err := pathSize(target)
			if err != nil {
				if len(args) == 1 {
					return err
				}
				fmt.Fprintf(os.Stderr, "warning: %s\n", err)
				continue
			}
			abs, err := filepath.Abs(target)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}
			rows = append(rows, row{formatutil.FormatLogicalSize(size), abs})
			total += size
		}

		if len(args) > 1 {
			rows = append(rows, row{formatutil.FormatLogicalSize(total), "total"})
		}

		// a tab misaligns once a size reaches the 8 column tab stop
		width := 0
		for _, r := range rows {
			width = max(width, len(r.size))
		}
		for _, r := range rows {
			fmt.Printf("%*s  %s\n", width, r.size, r.label)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sizeCmd)
}
