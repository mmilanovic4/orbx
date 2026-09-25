package cmd

import (
	"fmt"
	"strings"

	"github.com/mmilanovic4/orbx/internal/encodingutil"

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

		entries = append(entries, envEntry{key: key, value: parseEnvValue(value)})
	}

	return entries
}

// parseEnvValue strips the quotes around a value, or the comment after an
// unquoted one: in `A=1 # note` the value is 1, while in `A=a#b` the # is
// part of the value.
func parseEnvValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}

	if q := value[0]; q == '"' || q == '\'' {
		for i := 1; i < len(value); i++ {
			if q == '"' && value[i] == '\\' {
				i++ // skip the escaped character
				continue
			}
			if value[i] == q {
				return value[1:i]
			}
		}
		return value // no closing quote, show it as written
	}

	for i := 1; i < len(value); i++ {
		if value[i] == '#' && (value[i-1] == ' ' || value[i-1] == '\t') {
			return strings.TrimSpace(value[:i])
		}
	}
	return value
}

func init() {
	envCmd.Flags().StringVarP(&envFile, "file", "f", "", "read input from file")
	rootCmd.AddCommand(envCmd)
}
