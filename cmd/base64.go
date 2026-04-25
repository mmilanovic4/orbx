package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var (
	base64File string
)

var base64Cmd = &cobra.Command{
	Use:     "base64 [encode|decode] [text]",
	Short:   "Encode or decode base64",
	GroupID: "util",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := args[0]
		var input string

		if len(args) > 1 {
			input = args[1]
		}

		data, err := encodingutil.GetInputData(input, base64File)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		switch mode {
		case "encode":
			fmt.Println(encodingutil.EncodeBase64(data))
		case "decode":
			decoded, err := encodingutil.DecodeBase64(string(data))
			if err != nil {
				return fmt.Errorf("invalid base64 input: %w", err)
			}
			fmt.Println(string(decoded))
		default:
			return fmt.Errorf("unknown mode %q: use encode or decode", mode)
		}

		return nil
	},
}

func init() {
	base64Cmd.Flags().StringVarP(&base64File, "file", "f", "", "read input from file")
	rootCmd.AddCommand(base64Cmd)
}
