package sysutil

import (
	"fmt"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}
