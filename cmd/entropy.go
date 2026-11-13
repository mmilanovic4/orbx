package cmd

import (
	"fmt"
	"strings"

	"github.com/mmilanovic4/orbx/internal/cryptoutil"
	"github.com/mmilanovic4/orbx/internal/encodingutil"

	"github.com/spf13/cobra"
)

var (
	entropyFile   string
	entropyBase64 bool
)

var entropyCmd = &cobra.Command{
	Use:     "entropy [input]",
	Short:   "Calculate Shannon entropy of input",
	GroupID: "util",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, entropyFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		// decoding only on request: plenty of plain text (password, test,
		// 1234) is valid base64 too and would be measured as garbage bytes
		if entropyBase64 {
			decoded, err := encodingutil.DecodeBase64(strings.TrimSpace(string(data)))
			if err != nil {
				return fmt.Errorf("invalid base64 input: %w", err)
			}
			data = decoded
		}

		entropy := cryptoutil.ShannonEntropy(data)
		fmt.Printf("%.6f bits/byte\n", entropy)
		return nil
	},
}

func init() {
	entropyCmd.Flags().StringVarP(&entropyFile, "file", "f", "", "read input from file")
	entropyCmd.Flags().BoolVar(&entropyBase64, "base64", false, "decode base64 input and measure the decoded bytes")
	// raw input is the default now, the flag stays so existing scripts keep working
	entropyCmd.Flags().Bool("raw", false, "")
	entropyCmd.Flags().MarkDeprecated("raw", "input is measured as raw bytes by default, use --base64 to decode it first")
	rootCmd.AddCommand(entropyCmd)
}
