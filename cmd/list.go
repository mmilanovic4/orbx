package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows all available commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("orbx commands:")

		for _, c := range cmd.Root().Commands() {
			fmt.Printf("  %-12s %s\n", c.Name(), c.Short)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
