package cmd

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	uuidCompact bool
	uuidUpper   bool
)

var uuidCmd = &cobra.Command{
	Use:     "uuid",
	Short:   "Generate a UUID v4",
	GroupID: "dev",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := uuid.New().String()
		if uuidCompact {
			id = strings.ReplaceAll(id, "-", "")
		}
		if uuidUpper {
			id = strings.ToUpper(id)
		}
		fmt.Println(id)
		return nil
	},
}

func init() {
	uuidCmd.Flags().BoolVar(&uuidCompact, "compact", false, "remove dashes from UUID")
	uuidCmd.Flags().BoolVar(&uuidUpper, "upper", false, "uppercase UUID")
	rootCmd.AddCommand(uuidCmd)
}
