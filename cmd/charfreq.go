package cmd

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var (
	charfreqASCII bool
	charfreqFile  string
)

var charfreqCmd = &cobra.Command{
	Use:     "charfreq [input]",
	Short:   "Show character frequency table of input",
	GroupID: "util",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, charfreqFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		freq := make(map[rune]int)
		total := 0
		for _, r := range string(data) {
			freq[r]++
			total++
		}

		type entry struct {
			char  rune
			count int
		}

		entries := make([]entry, 0, len(freq))
		for r, count := range freq {
			entries = append(entries, entry{r, count})
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].count > entries[j].count
		})

		maxCount := entries[0].count

		if charfreqASCII {
			fmt.Printf("%-8s %-8s %-8s %s\n", "ASCII", "Count", "%", "Bar")
		} else {
			fmt.Printf("%-8s %-8s %-8s %s\n", "Char", "Count", "%", "Bar")
		}
		fmt.Println(strings.Repeat("─", 48))

		for _, e := range entries {
			pct := float64(e.count) / float64(total) * 100
			barWidth := int(float64(e.count) / float64(maxCount) * 20)
			bar := strings.Repeat("█", barWidth)

			var label string
			if charfreqASCII {
				label = fmt.Sprintf("%d", e.char)
			} else {
				if unicode.IsControl(e.char) {
					label = fmt.Sprintf("\\x%02x", e.char)
				} else {
					label = string(e.char)
				}
			}

			fmt.Printf("%-8s %-8d %-8.2f %s\n", label, e.count, pct, bar)
		}

		fmt.Printf("\n%d unique characters, %d total\n", len(entries), total)

		return nil
	},
}

func init() {
	charfreqCmd.Flags().BoolVar(&charfreqASCII, "ascii", false, "show ASCII values instead of characters")
	charfreqCmd.Flags().StringVarP(&charfreqFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(charfreqCmd)
}
