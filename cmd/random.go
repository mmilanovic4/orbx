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
	Run: func(cmd *cobra.Command, args []string) {
		length, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
			fmt.Println("length is not valid")
			return
		}

		if length <= 0 || length > MAX_LENGTH {
			fmt.Printf("length not valid, should be between 0 and %d\n", MAX_LENGTH)
			return
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
				fmt.Println("failed to generate random string:", err)
				return
			}
			result[i] = charset[n.Int64()]
		}

		fmt.Println(string(result))
	},
}

func init() {
	randomCmd.Flags().StringVarP(&randomCharset, "charset", "c", "alphanum", "character set: alpha, alphanum, hex, symbols")
	rootCmd.AddCommand(randomCmd)
}
