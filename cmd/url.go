package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	urlTyped bool
)

func parseTypedValue(v string) any {
	if v == "" {
		return true
	}

	if v == "null" || v == "nil" || v == "undefined" {
		return nil
	}

	switch strings.ToLower(v) {
	case "true":
		return true
	case "false":
		return false
	}

	// numbers are typed only if that keeps them exactly as written, so
	// values like 007, 1.50 or integers beyond int64 stay strings
	if i, err := strconv.ParseInt(v, 10, 64); err == nil && strconv.FormatInt(i, 10) == v {
		return i
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil && strconv.FormatFloat(f, 'f', -1, 64) == v {
		return f
	}

	return v
}

// parseURL parses raw before anything is decoded, so escaped separators in
// values (%26, %3D, %23, %2B) stay part of the value. A URL that was
// percent-encoded as a whole, e.g. copied from a redirect parameter, is
// decoded once first. The returned string is the decoded URL for display.
func parseURL(raw string) (*url.URL, string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, "", fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme == "" && u.Host == "" {
		if unescaped, err := url.QueryUnescape(raw); err == nil && unescaped != raw {
			if inner, err := url.Parse(unescaped); err == nil && inner.Scheme != "" && inner.Host != "" {
				u, raw = inner, unescaped
			}
		}
	}

	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return nil, "", fmt.Errorf("invalid URL encoding: %w", err)
	}

	return u, decoded, nil
}

var urlCmd = &cobra.Command{
	Use:     "url [url]",
	Short:   "Decode and parse a URL",
	GroupID: "dev",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		u, decoded, err := parseURL(args[0])
		if err != nil {
			return err
		}

		fmt.Println("Decoded URL:", decoded)
		fmt.Println("Scheme:", u.Scheme)
		fmt.Println("Host:", u.Host)
		fmt.Println("Path:", u.Path)
		if u.Fragment != "" {
			fmt.Println("Fragment:", u.Fragment)
		}

		queryMap := make(map[string]any)

		for key, values := range u.Query() {
			if len(values) == 1 {
				val := values[0]
				if urlTyped {
					queryMap[key] = parseTypedValue(val)
				} else {
					queryMap[key] = val
				}
			} else {
				if urlTyped {
					typed := make([]any, len(values))
					for i, v := range values {
						typed[i] = parseTypedValue(v)
					}
					queryMap[key] = typed
				} else {
					queryMap[key] = values
				}
			}
		}

		// decoded values are shown as they are, not with & < > escaped
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(queryMap); err != nil {
			return fmt.Errorf("failed to encode query: %w", err)
		}

		fmt.Println("Query params:")
		fmt.Print(buf.String())

		return nil
	},
}

func init() {
	urlCmd.Flags().BoolVar(&urlTyped, "typed", false, "parse query params into typed values")
	rootCmd.AddCommand(urlCmd)
}
