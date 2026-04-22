package cmd

import (
	"fmt"
	"orbx/internal/cryptoutil"
	"orbx/internal/encodingutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	entropyFile string
	entropyRaw  bool
)

var entropyCmd = &cobra.Command{
	Use:     "entropy [input]",
	Short:   "Calculate Shannon entropy of input",
	GroupID: "util",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, entropyFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if !entropyRaw {
			decoded, err := encodingutil.DecodeBase64(strings.TrimSpace(string(data)))
			if err == nil {
				data = decoded
			}
		}

		entropy := cryptoutil.ShannonEntropy(data)
		fmt.Printf("%.6f bits/byte\n", entropy)
	},
}

func init() {
	entropyCmd.Flags().StringVarP(&entropyFile, "file", "f", "", "read input from file")
	entropyCmd.Flags().BoolVar(&entropyRaw, "raw", false, "treat input as raw bytes, skip base64 decoding")
	rootCmd.AddCommand(entropyCmd)
}
