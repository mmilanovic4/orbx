package cmd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"orbx/internal/encodingutil"
	"orbx/internal/sysutil"
	"os"

	"github.com/spf13/cobra"
)

var (
	compressFile   string
	compressOut    string
	compressDecode bool
	compressRaw    bool
)

var compressCmd = &cobra.Command{
	Use:     "compress [input]",
	Short:   "Compress or decompress input using gzip",
	GroupID: "util",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var rawInput string
		if len(args) > 0 {
			rawInput = args[0]
		}

		data, err := encodingutil.GetInputData(rawInput, compressFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		var result []byte

		if compressDecode {
			result, err = gunzip(data, compressRaw)
			if err != nil {
				return fmt.Errorf("failed to decompress: %w", err)
			}
			fmt.Println(string(result))
			if compressOut != "" {
				if err := sysutil.WriteFile(compressOut, result); err != nil {
					return fmt.Errorf("failed to write output file: %w", err)
				}
			}
		} else {
			result, err = gzipCompress(data)
			if err != nil {
				return fmt.Errorf("failed to compress: %w", err)
			}
			if compressRaw {
				os.Stdout.Write(result)
				if compressOut != "" {
					if err := sysutil.WriteFile(compressOut, result); err != nil {
						return fmt.Errorf("failed to write output file: %w", err)
					}
				}
			} else {
				encoded := encodingutil.EncodeBase64(result)
				fmt.Println(encoded)
				if compressOut != "" {
					if err := sysutil.WriteFile(compressOut, []byte(encoded)); err != nil {
						return fmt.Errorf("failed to write output file: %w", err)
					}
				}
			}
		}

		return nil
	},
}

func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gunzip(data []byte, raw bool) ([]byte, error) {
	if !raw {
		decoded, err := encodingutil.DecodeBase64(string(data))
		if err == nil {
			data = decoded
		}
	}

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

func init() {
	compressCmd.Flags().StringVarP(&compressFile, "file", "f", "", "input file (optional)")
	compressCmd.Flags().StringVarP(&compressOut, "out", "o", "", "output file (optional)")
	compressCmd.Flags().BoolVarP(&compressDecode, "decode", "d", false, "decompress instead of compress")
	compressCmd.Flags().BoolVar(&compressRaw, "raw", false, "write raw gzip bytes instead of base64 (compatible with 7-Zip, gunzip, etc.)")
	rootCmd.AddCommand(compressCmd)
}
