package cmd

import (
	"encoding/json"
	"fmt"
	"orbx/internal/encodingutil"
	"strings"

	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:     "jwt [token]",
	Short:   "Decode a JWT token (header and payload, no verification)",
	GroupID: "dev",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			fmt.Println("Invalid JWT: expected 3 parts separated by '.'")
			return
		}

		header, err := decodeJWTPart(parts[0])
		if err != nil {
			fmt.Println("Failed to decode header:", err)
			return
		}

		payload, err := decodeJWTPart(parts[1])
		if err != nil {
			fmt.Println("Failed to decode payload:", err)
			return
		}

		fmt.Println("=== Header ===")
		fmt.Println(header)
		fmt.Println("=== Payload ===")
		fmt.Println(payload)
	},
}

func decodeJWTPart(part string) (string, error) {
	// JWT uses base64url (no padding), add padding if needed
	switch len(part) % 4 {
	case 2:
		part += "=="
	case 3:
		part += "="
	}

	data, err := encodingutil.DecodeBase64(part)
	if err != nil {
		return "", err
	}

	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		return "", err
	}

	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}

	return string(pretty), nil
}

func init() {
	rootCmd.AddCommand(jwtCmd)
}
