package sysutil

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteOutputIsByteExact(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"text without newline", []byte("abc")},
		{"text with newline", []byte("abc\n")},
		{"binary", []byte{0x1f, 0x8b, 0x00, 0xff}},
		{"empty", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "out")
			f, err := os.Create(path)
			if err != nil {
				t.Fatal(err)
			}

			if err := writeOutput(f, tt.data); err != nil {
				t.Fatalf("writeOutput() error = %v", err)
			}
			f.Close()

			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, tt.data) {
				t.Errorf("wrote %q, want %q", got, tt.data)
			}
		})
	}
}
