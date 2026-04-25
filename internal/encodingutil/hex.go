package encodingutil

import (
	"encoding/hex"
	"fmt"
)

func EncodeHex(data []byte) string {
	return hex.EncodeToString(data)
}

func DecodeHex(input string) ([]byte, error) {
	data, err := hex.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("invalid hex input: %w", err)
	}
	return data, nil
}
