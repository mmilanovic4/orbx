package encodingutil

import (
	"encoding/hex"
	"strings"
)

func EncodeHex(data []byte) string {
	return hex.EncodeToString(data)
}

// DecodeHex ignores whitespace, so input with a trailing newline (echo) or
// split over lines (xxd -p) decodes as well.
func DecodeHex(input string) ([]byte, error) {
	return hex.DecodeString(strings.Join(strings.Fields(input), ""))
}
