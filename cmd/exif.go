package cmd

import (
	"fmt"

	"github.com/mmilanovic4/orbx/internal/imageutil"
	"github.com/mmilanovic4/orbx/internal/sysutil"

	"github.com/spf13/cobra"
)

var exifAll bool

var exifCmd = &cobra.Command{
	Use:     "exif [file]",
	Short:   "Show EXIF metadata tags of an image",
	GroupID: "util",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := sysutil.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		tags, err := imageutil.ReadEXIF(data, exifAll)
		if err != nil {
			return err
		}

		width := 0
		for _, t := range tags {
			if len(t.Name) > width {
				width = len(t.Name)
			}
		}

		for _, t := range tags {
			fmt.Printf("%-*s %s\n", width+1, t.Name+":", t.Value)
		}

		return nil
	},
}

func init() {
	exifCmd.Flags().BoolVarP(&exifAll, "all", "a", false, "include tags without a known name")
	rootCmd.AddCommand(exifCmd)
}
