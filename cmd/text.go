package cmd

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var (
	substring string
)

var textCmd = &cobra.Command{
	Use:     "text [operation] [input]",
	Short:   "String utilities",
	GroupID: "util",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		op := args[0]
		input := args[1]

		switch op {
		case "upper":
			fmt.Println(strings.ToUpper(input))
		case "lower":
			fmt.Println(strings.ToLower(input))
		case "title":
			fmt.Println(titleCase(input))
		case "trim":
			fmt.Println(strings.TrimSpace(input))
		case "reverse":
			runes := []rune(input)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			fmt.Println(string(runes))
		case "slug":
			s := strings.ToLower(input)
			var b strings.Builder
			prevDash := false
			for _, r := range s {
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					b.WriteRune(r)
					prevDash = false
				} else if !prevDash {
					b.WriteRune('-')
					prevDash = true
				}
			}
			fmt.Println(strings.Trim(b.String(), "-"))
		case "count":
			fmt.Println(len([]rune(input)))
		case "words":
			fmt.Println(len(strings.Fields(input)))
		case "contains":
			if !cmd.Flags().Changed("substring") {
				return fmt.Errorf("contains requires --substring (-s)")
			}
			fmt.Println(strings.Contains(input, substring))
		default:
			return fmt.Errorf("unknown operation %q — use: upper, lower, title, trim, reverse, slug, count, words, contains [-s substring]", op)
		}

		return nil
	},
}

// titleCase upper-cases the first letter of every whitespace-separated word
// and lower-cases the rest, e.g. "hello WORLD" becomes "Hello World".
// strings.ToTitle is not an option: it maps every letter to title case.
func titleCase(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	wordStart := true
	for _, r := range s {
		switch {
		case unicode.IsSpace(r):
			wordStart = true
		case wordStart:
			r = unicode.ToTitle(r)
			wordStart = false
		default:
			r = unicode.ToLower(r)
		}
		b.WriteRune(r)
	}

	return b.String()
}

func init() {
	textCmd.Flags().StringVarP(&substring, "substring", "s", "", "Substring to search for")
	rootCmd.AddCommand(textCmd)
}
