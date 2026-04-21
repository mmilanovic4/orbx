package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var htmlCmd = &cobra.Command{
	Use:     "html [encode|decode] [text]",
	Short:   "Encode or decode HTML entities",
	GroupID: "dev",
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]

		var input string
		if len(args) > 1 {
			input = args[1]
		}

		data, err := encodingutil.GetInputData(input, "")
		if err != nil {
			fmt.Println(err)
			return
		}

		text := string(data)

		switch action {
		case "encode":
			fmt.Println(encodingutil.EncodeHTML(text))
		case "decode":
			fmt.Println(encodingutil.DecodeHTML(text))
		default:
			fmt.Println("Use encode or decode")
		}
	},
}

func init() {
	rootCmd.AddCommand(htmlCmd)
}
