package cmd

import (
	"fmt"

	"github.com/mmilanovic4/orbx/internal/encodingutil"
	"github.com/mmilanovic4/orbx/internal/sysutil"

	"github.com/spf13/cobra"
)

var (
	hexFile string
)

var hexCmd = &cobra.Command{
	Use:     "hex [encode|decode] [text]",
	Short:   "Encode or decode hex",
	GroupID: "util",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := args[0]
		var input string

		if len(args) > 1 {
			input = args[1]
		}

		data, err := encodingutil.GetInputData(input, hexFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		switch mode {
		case "encode":
			fmt.Println(encodingutil.EncodeHex(data))
		case "decode":
			decoded, err := encodingutil.DecodeHex(string(data))
			if err != nil {
				return fmt.Errorf("invalid hex input: %w", err)
			}
			if err := sysutil.WriteStdout(decoded); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		default:
			return fmt.Errorf("unknown mode %q: use encode or decode", mode)
		}

		return nil
	},
}

func init() {
	hexCmd.Flags().StringVarP(&hexFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(hexCmd)
}
