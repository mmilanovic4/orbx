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
		fmt.Println("clearclip			clear system clipboard")
		fmt.Println("list						show available commands")
		fmt.Println("version				show version")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
