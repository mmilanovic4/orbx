package cmd

import (
	"fmt"
	"orbx/internal/encodingutil"
	"strings"

	"github.com/spf13/cobra"
)

var (
	envFile string
)

type envEntry struct {
	key   string
	value string
}

var envCmd = &cobra.Command{
	Use:     "env [input]",
	Short:   "Pretty print env file content",
	GroupID: "dev",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		}

		data, err := encodingutil.GetInputData(input, envFile)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		entries := parseEnv(string(data))
		if len(entries) == 0 {
			return fmt.Errorf("no variables found in input")
		}

		width := 0
		for _, e := range entries {
			if len(e.key) > width {
				width = len(e.key)
			}
		}

		for _, e := range entries {
			fmt.Printf("%-*s  %s\n", width, e.key, e.value)
		}

		return nil
	},
}

func parseEnv(content string) []envEntry {
	var entries []envEntry

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Strip surrounding quotes
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		entries = append(entries, envEntry{key: key, value: value})
	}

	return entries
}

func init() {
	envCmd.Flags().StringVarP(&envFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(envCmd)
}
