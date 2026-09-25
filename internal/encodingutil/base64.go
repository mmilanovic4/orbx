package encodingutil

import (
	"encoding/base64"
	"strings"
)

func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeBase64(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

// DecodeBase64URL decodes the URL-safe alphabet used by JWTs, with or
// without padding. Input in the standard alphabet is accepted as well.
func DecodeBase64URL(input string) ([]byte, error) {
	input = strings.TrimRight(input, "=")
	input = strings.NewReplacer("+", "-", "/", "_").Replace(input)
	return base64.RawURLEncoding.DecodeString(input)
}
