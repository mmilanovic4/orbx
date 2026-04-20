package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var hexFile string

var hexCmd = &cobra.Command{
	Use:     "hex [encode|decode] [text]",
	Short:   "Encode or decode hex",
	GroupID: "dev",
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		mode := args[0]
		var input string

		if len(args) > 1 {
			input = args[1]
		}

		data, err := encodingutil.GetInputData(input, hexFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch mode {
		case "encode":
			fmt.Println(encodingutil.EncodeHex(data))
		case "decode":
			decoded, err := encodingutil.DecodeHex(string(data))
			if err != nil {
				fmt.Println("Invalid base64 input")
				return
			}
			fmt.Println(string(decoded))
		default:
			fmt.Println("Use encode or decode")
		}
	},
}

func init() {
	hexCmd.Flags().StringVarP(&hexFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(hexCmd)
}
