package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"orbx/internal/encodingutil"
	"strings"

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
			var obj any
			if err := json.Unmarshal(data, &obj); err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}
			pretty, err := json.MarshalIndent(obj, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(pretty))
			return nil
		}

		if strings.HasPrefix(trimmed, "<") {
			var buf strings.Builder
			decoder := xml.NewDecoder(strings.NewReader(trimmed))
			encoder := xml.NewEncoder(&buf)
			encoder.Indent("", "  ")
			for {
				token, err := decoder.Token()
				if err != nil {
					break
				}
				if err := encoder.EncodeToken(token); err != nil {
					return fmt.Errorf("invalid XML: %w", err)
				}
			}
			encoder.Flush()
			fmt.Println(buf.String())
			return nil
		}

		return fmt.Errorf("unsupported format: input must be JSON or XML")
	},
}

func init() {
	prettyCmd.Flags().StringVarP(&prettyFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(prettyCmd)
}
