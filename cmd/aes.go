package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"orbx/internal/cryptoutil"
	"orbx/internal/encodingutil"
	"orbx/internal/sysutil"

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
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		keyFile, _ := cmd.Flags().GetString("key")
		outFile, _ := cmd.Flags().GetString("out")

		key, err := cryptoutil.ReadKey(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read key: %w", err)
		}

		plainText, err := encodingutil.GetInputData("", input)
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
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		keyFile, _ := cmd.Flags().GetString("key")
		outFile, _ := cmd.Flags().GetString("out")

		key, err := cryptoutil.ReadKey(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read key: %w", err)
		}

		cipherTextEncoded, err := encodingutil.GetInputData("", input)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		cipherText, err := encodingutil.DecodeBase64(string(cipherTextEncoded))
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
	Use:   "key",
	Short: "Generate a new AES-256 key",
	RunE: func(cmd *cobra.Command, args []string) error {
		outFile, _ := cmd.Flags().GetString("out")

		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		keyEncoded := base64.StdEncoding.EncodeToString(key)
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
	aesEncryptCmd.Flags().String("out", "", "Output file (optional)")
	aesEncryptCmd.MarkFlagRequired("key")

	aesDecryptCmd.Flags().String("key", "", "Path to key file (required)")
	aesDecryptCmd.Flags().String("out", "", "Output file (optional)")
	aesDecryptCmd.MarkFlagRequired("key")

	aesKeyCmd.Flags().String("out", "", "Output file (optional)")
}
