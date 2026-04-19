package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var hashFile string

var hashCmd = &cobra.Command{
	Use:     "hash [algorithm] [text]",
	Short:   "Generate hash of a string",
	GroupID: "dev",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		algo := args[0]
		var input string

		if len(args) > 1 {
			input = args[1]
		}

		data, err := encodingutil.GetInputData(input, base64File)
		if err != nil {
			fmt.Println(err)
			return
		}

		hash, err := encodingutil.Hash(algo, data)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(encodingutil.EncodeHex(hash))
	},
}

func init() {
	hashCmd.Flags().StringVarP(&base64File, "file", "f", "", "read input from file")
	rootCmd.AddCommand(hashCmd)
}
