package netutil

import "strings"

// NormalizeURL adds https:// to input that has no http or https scheme. The
// scheme is matched case-insensitively, as URL schemes are.
func NormalizeURL(input string) string {
	lower := strings.ToLower(input)
	if !strings.HasPrefix(lower, "http://") && !strings.HasPrefix(lower, "https://") {
		return "https://" + input
	}
	return input
}
