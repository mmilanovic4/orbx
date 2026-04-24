package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var htmlFile string

var htmlCmd = &cobra.Command{
	Use:     "html [encode|decode] [text]",
	Short:   "Encode or decode HTML entities",
	GroupID: "dev",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]

		var input string
		if len(args) > 1 {
			input = args[1]
		}

		data, err := encodingutil.GetInputData(input, htmlFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		text := string(data)

		switch action {
		case "encode":
			fmt.Println(encodingutil.EncodeHTML(text))
		case "decode":
			fmt.Println(encodingutil.DecodeHTML(text))
		default:
			return fmt.Errorf("unknown mode %q: use encode or decode", action)
		}

		return nil
	},
}

func init() {
	htmlCmd.Flags().StringVarP(&htmlFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(htmlCmd)
}
