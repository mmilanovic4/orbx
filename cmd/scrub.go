package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmilanovic4/orbx/internal/formatutil"
	"github.com/mmilanovic4/orbx/internal/imageutil"
	"github.com/mmilanovic4/orbx/internal/sysutil"

	"github.com/spf13/cobra"
)

var (
	scrubOut      string
	scrubForce    bool
	scrubStripICC bool
)

var scrubCmd = &cobra.Command{
	Use:     "scrub [file]",
	Short:   "Remove EXIF and other metadata from an image (lossless)",
	GroupID: "util",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := sysutil.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		cleaned, removed, err := imageutil.Scrub(data, scrubStripICC)
		if err != nil {
			return err
		}

		if len(removed) == 0 {
			fmt.Println("No metadata found, nothing to remove.")
			return nil
		}

		out := scrubOut
		if out == "" {
			out = scrubOutputPath(args[0])
		}

		if !scrubForce {
			if _, err := os.Stat(out); err == nil {
				return fmt.Errorf("%s already exists, use --force to overwrite", out)
			}
		}

		if err := sysutil.WriteFile(out, cleaned); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		for _, r := range removed {
			fmt.Printf("✓ Removed %s (%s)\n", r.Name, formatutil.FormatLogicalSize(int64(r.Bytes)))
		}

		fmt.Printf("Saved to %s (%s → %s)\n", out,
			formatutil.FormatLogicalSize(int64(len(data))),
			formatutil.FormatLogicalSize(int64(len(cleaned))))

		return nil
	},
}

func scrubOutputPath(input string) string {
	ext := filepath.Ext(input)
	return strings.TrimSuffix(input, ext) + "-clean" + ext
}

func init() {
	scrubCmd.Flags().StringVarP(&scrubOut, "out", "o", "", "output file path (default: <name>-clean.<ext>)")
	scrubCmd.Flags().BoolVarP(&scrubForce, "force", "f", false, "overwrite the output file if it exists")
	scrubCmd.Flags().BoolVar(&scrubStripICC, "strip-icc", false, "also remove the ICC color profile")
	rootCmd.AddCommand(scrubCmd)
}
