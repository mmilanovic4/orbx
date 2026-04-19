package netutil

import "strings"

func NormalizeURL(input string) string {
	if !strings.HasPrefix(input, "http://") &&
		!strings.HasPrefix(input, "https://") {
		return "https://" + input
	}
	return input
}
