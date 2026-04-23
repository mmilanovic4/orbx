package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"orbx/internal/encodingutil"
	"strings"

	"github.com/spf13/cobra"
)

var prettyFile string

var prettyCmd = &cobra.Command{
	Use:     "prettyprint [input]",
	Short:   "Format and pretty print JSON or XML",
	GroupID: "dev",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, prettyFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		trimmed := strings.TrimSpace(string(data))

		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			var obj any
			if err := json.Unmarshal(data, &obj); err != nil {
				fmt.Println("Invalid JSON:", err)
				return
			}
			pretty, err := json.MarshalIndent(obj, "", "  ")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(pretty))
			return
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
					fmt.Println("Invalid XML:", err)
					return
				}
			}
			encoder.Flush()
			fmt.Println(buf.String())
			return
		}

		fmt.Println("Unsupported format: input must be JSON or XML")
	},
}

func init() {
	prettyCmd.Flags().StringVarP(&prettyFile, "file", "f", "", "read input from file")
	prettyCmd.Aliases = []string{"pp"}
	rootCmd.AddCommand(prettyCmd)
}
