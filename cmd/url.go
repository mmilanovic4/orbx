package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var urlTyped bool

func parseTypedValue(v string) any {
	if v == "" {
		return true
	}

	// null/undefined
	if v == "null" || v == "nil" || v == "undefined" {
		return nil
	}

	// bool (case-insensitive)
	switch strings.ToLower(v) {
	case "true":
		return true
	case "false":
		return false
	}

	// float64
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		if strconv.FormatFloat(f, 'g', -1, 64) == v {
			return f
		}
		return v // big numbers are kept as strings
	}

	return v
}

var urlCmd = &cobra.Command{
	Use:   "url [url]",
	Short: "Decode and parse a URL",
	Run: func(cmd *cobra.Command, args []string) {
		raw := args[0]

		decoded, err := url.QueryUnescape(raw)
		if err != nil {
			fmt.Println("Invalid URL encoding:", err)
			return
		}

		fmt.Println("Decoded URL:", decoded)

		u, err := url.Parse(decoded)
		if err != nil {
			fmt.Println("Invalid URL:", err)
			return
		}

		fmt.Println("Scheme:", u.Scheme)
		fmt.Println("Host:", u.Host)
		fmt.Println("Path:", u.Path)

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

		jsonBytes, err := json.MarshalIndent(queryMap, "", "  ")
		if err != nil {
			fmt.Println("Failed to encode query:", err)
			return
		}

		fmt.Println("Query params:")
		fmt.Println(string(jsonBytes))

	},
}

func init() {
	urlCmd.Flags().BoolVar(&urlTyped, "typed", false, "parse query params into typed values")
	rootCmd.AddCommand(urlCmd)
}
