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
	GroupID: "util",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		algo := args[0]
		var input string

		if len(args) > 1 {
			input = args[1]
		}

		data, err := encodingutil.GetInputData(input, hashFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		hash, err := encodingutil.Hash(algo, data)
		if err != nil {
			return fmt.Errorf("failed to generate hash: %w", err)
		}

		fmt.Println(encodingutil.EncodeHex(hash))
		return nil
	},
}

func init() {
	hashCmd.Flags().StringVarP(&hashFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(hashCmd)
}
