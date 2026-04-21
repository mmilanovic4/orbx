package encodingutil

import (
	"fmt"
	"orbx/internal/sysutil"
)

func GetInputData(text string, file string) ([]byte, error) {
	if file != "" {
		data, err := sysutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		return data, nil
	}

	if text != "" {
		return []byte(text), nil
	}

	return nil, fmt.Errorf("no input provided")
}
