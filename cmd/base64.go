package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var base64File string

var base64Cmd = &cobra.Command{
	Use:     "base64 [encode|decode] [text]",
	Short:   "Encode or decode base64",
	GroupID: "dev",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mode := args[0]
		var input string

		if len(args) > 1 {
			input = args[1]
		}

		data, err := encodingutil.GetInputData(input, base64File)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch mode {
		case "encode":
			fmt.Println(encodingutil.EncodeBase64(data))
		case "decode":
			decoded, err := encodingutil.DecodeBase64(string(data))
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
	base64Cmd.Flags().StringVarP(&base64File, "file", "f", "", "read input from file")
	rootCmd.AddCommand(base64Cmd)
}
