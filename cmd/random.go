package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	charsetAlpha    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetAlphanum = charsetAlpha + "0123456789"
	charsetHex      = "0123456789abcdef"
	charsetSymbols  = charsetAlphanum + "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

var randomCharset string

const MAX_LENGTH int = 10000

var randomCmd = &cobra.Command{
	Use:     "random [length]",
	Short:   "Generate a cryptographically secure random string",
	GroupID: "util",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		length, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid length %q: must be an integer", args[0])
		}

		if length <= 0 || length > MAX_LENGTH {
			return fmt.Errorf("length must be between 1 and %d", MAX_LENGTH)
		}

		var charset string
		switch strings.ToLower(randomCharset) {
		case "alpha":
			charset = charsetAlpha
		case "alphanum":
			charset = charsetAlphanum
		case "hex":
			charset = charsetHex
		case "symbols":
			charset = charsetSymbols
		default:
			charset = charsetAlphanum
		}

		result := make([]byte, length)
		for i := range result {
			n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return fmt.Errorf("failed to generate random string: %w", err)
			}
			result[i] = charset[n.Int64()]
		}

		fmt.Println(string(result))
		return nil
	},
}

func init() {
	randomCmd.Flags().StringVarP(&randomCharset, "charset", "c", "alphanum", "character set: alpha, alphanum, hex, symbols")
	rootCmd.AddCommand(randomCmd)
}
