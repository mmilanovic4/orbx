package sysutil

import (
	"bytes"
	"fmt"
	"io"
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

// ReplaceFile is WriteFile for outputs that may overwrite a file the user
// cares about, even the input itself: the data goes to a temporary file in
// the same directory that is then renamed over path, so an interrupted
// write never leaves path half written.
func ReplaceFile(path string, data []byte) error {
	return ReplaceFileFrom(path, bytes.NewReader(data))
}

// ReplaceFileFrom is ReplaceFile for data that is streamed, such as a
// download: it is copied from r without being held in memory.
func ReplaceFileFrom(path string, r io.Reader) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	tmp := f.Name()

	err = f.Chmod(0644)
	if err == nil {
		_, err = io.Copy(f, r)
	}
	if err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		err = os.Rename(tmp, path)
	}
	if err != nil {
		os.Remove(tmp)
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}
