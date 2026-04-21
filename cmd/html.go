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
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]
		text := args[1]

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
