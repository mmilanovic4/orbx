package cmd

import (
	"crypto/rand"
	"fmt"
	"orbx/internal/cryptoutil"
	"orbx/internal/encodingutil"
	"orbx/internal/sysutil"
	"strings"

	"github.com/spf13/cobra"
)

var aesCmd = &cobra.Command{
	Use:     "aes",
	Short:   "AES-GCM encryption utilities",
	GroupID: "util",
}

var aesEncryptCmd = &cobra.Command{
	Use:   "encrypt [input]",
	Short: "Encrypt input using AES-GCM",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var rawInput string
		if len(args) > 0 {
			rawInput = args[0]
		}

		keyFile, _ := cmd.Flags().GetString("key")
		inputFile, _ := cmd.Flags().GetString("file")
		outFile, _ := cmd.Flags().GetString("out")

		key, err := cryptoutil.ReadKey(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read key: %w", err)
		}

		plainText, err := encodingutil.GetInputData(rawInput, inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		cipherText, err := cryptoutil.Encrypt(plainText, key)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		cipherTextEncoded := encodingutil.EncodeBase64(cipherText)
		fmt.Println(cipherTextEncoded)

		if outFile != "" {
			if err := sysutil.WriteFile(outFile, []byte(cipherTextEncoded)); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
		}

		return nil
	},
}

var aesDecryptCmd = &cobra.Command{
	Use:   "decrypt [input]",
	Short: "Decrypt input using AES-GCM",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var rawInput string
		if len(args) > 0 {
			rawInput = args[0]
		}

		keyFile, _ := cmd.Flags().GetString("key")
		inputFile, _ := cmd.Flags().GetString("file")
		outFile, _ := cmd.Flags().GetString("out")

		key, err := cryptoutil.ReadKey(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read key: %w", err)
		}

		rawBytes, err := encodingutil.GetInputData(rawInput, inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		cipherText, err := encodingutil.DecodeBase64(strings.TrimSpace(string(rawBytes)))
		if err != nil {
			return fmt.Errorf("failed to decode base64 input: %w", err)
		}

		plainText, err := cryptoutil.Decrypt(cipherText, key)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		fmt.Println(string(plainText))

		if outFile != "" {
			if err := sysutil.WriteFile(outFile, plainText); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
		}

		return nil
	},
}

var aesKeyCmd = &cobra.Command{
	Use:   "key [size]",
	Short: "Generate a new AES key (size: 16, 24, 32 — default 32)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outFile, _ := cmd.Flags().GetString("out")

		size := 32
		if len(args) == 1 {
			switch args[0] {
			case "16":
				size = 16
			case "24":
				size = 24
			case "32":
				size = 32
			default:
				size = 32
			}
		}

		key := make([]byte, size)
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		keyEncoded := encodingutil.EncodeBase64(key)
		fmt.Println(keyEncoded)

		if outFile != "" {
			if err := sysutil.WriteFile(outFile, []byte(keyEncoded)); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(aesCmd)

	aesCmd.AddCommand(aesEncryptCmd)
	aesCmd.AddCommand(aesDecryptCmd)
	aesCmd.AddCommand(aesKeyCmd)

	aesEncryptCmd.Flags().String("key", "", "Path to key file (required)")
	aesEncryptCmd.Flags().String("file", "", "Input file (optional)")
	aesEncryptCmd.Flags().String("out", "", "Output file (optional)")
	aesEncryptCmd.MarkFlagRequired("key")

	aesDecryptCmd.Flags().String("key", "", "Path to key file (required)")
	aesDecryptCmd.Flags().String("file", "", "Input file (optional)")
	aesDecryptCmd.Flags().String("out", "", "Output file (optional)")
	aesDecryptCmd.MarkFlagRequired("key")

	aesKeyCmd.Flags().String("out", "", "Output file (optional)")
}
