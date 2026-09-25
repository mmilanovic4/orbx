package cmd

import (
	"fmt"
	"strings"

	"github.com/mmilanovic4/orbx/internal/encodingutil"
	"github.com/mmilanovic4/orbx/internal/formatutil"

	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:     "jwt [token]",
	Short:   "Decode a JWT token (header and payload, no verification)",
	GroupID: "dev",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			return fmt.Errorf("invalid JWT: expected 3 parts separated by '.'")
		}

		header, err := decodeJWTPart(parts[0])
		if err != nil {
			return fmt.Errorf("failed to decode header: %w", err)
		}

		payload, err := decodeJWTPart(parts[1])
		if err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}

		fmt.Println("=== Header ===")
		fmt.Println(header)
		fmt.Println("=== Payload ===")
		fmt.Println(payload)

		return nil
	},
}

func decodeJWTPart(part string) (string, error) {
	data, err := encodingutil.DecodeBase64URL(part)
	if err != nil {
		return "", err
	}

	return formatutil.IndentJSON(data)
}

func init() {
	rootCmd.AddCommand(jwtCmd)
}
