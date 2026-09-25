package cmd

import (
	"fmt"
	"strings"

	"github.com/mmilanovic4/orbx/internal/encodingutil"
	"github.com/mmilanovic4/orbx/internal/formatutil"

	"github.com/spf13/cobra"
)

var (
	prettyFile string
)

var prettyCmd = &cobra.Command{
	Use:     "prettyprint [input]",
	Short:   "Format and pretty print JSON or XML",
	GroupID: "dev",
	Aliases: []string{"pp"},
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, prettyFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		trimmed := strings.TrimSpace(string(data))

		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			pretty, err := formatutil.IndentJSON(data)
			if err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}
			fmt.Println(pretty)
			return nil
		}

		if strings.HasPrefix(trimmed, "<") {
			pretty, err := formatutil.IndentXML(data)
			if err != nil {
				return fmt.Errorf("invalid XML: %w", err)
			}
			fmt.Println(pretty)
			return nil
		}

		return fmt.Errorf("unsupported format: input must be JSON or XML")
	},
}

func init() {
	prettyCmd.Flags().StringVarP(&prettyFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(prettyCmd)
}
