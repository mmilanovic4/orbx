package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var wcFile string

var wcCmd = &cobra.Command{
	Use:     "wc [input]",
	Short:   "Count lines, words and characters",
	GroupID: "util",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, wcFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		text := string(data)

		lines := strings.Count(text, "\n")
		words := len(strings.FieldsFunc(text, unicode.IsSpace))
		chars := len([]rune(text))
		bytes_ := len(data)

		fmt.Printf("lines: %d\n", lines)
		fmt.Printf("words: %d\n", words)
		fmt.Printf("chars: %d\n", chars)
		fmt.Printf("bytes: %d\n", bytes_)

		return nil
	},
}

func init() {
	wcCmd.Flags().StringVarP(&wcFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(wcCmd)
}
