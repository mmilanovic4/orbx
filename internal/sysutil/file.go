package sysutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}

// WritePrivateFile is WriteFile for secrets such as keys and decrypted
// data: the file ends up readable by its owner only.
func WritePrivateFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// an existing file keeps its old mode on open, so tighten it before
	// anything is written
	if err := f.Chmod(0600); err != nil {
		f.Close()
		return fmt.Errorf("failed to save file: %w", err)
	}

	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("failed to save file: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}
