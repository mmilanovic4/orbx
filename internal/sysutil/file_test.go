package sysutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWritePrivateFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permissions")
	}

	tests := []struct {
		name     string
		existing os.FileMode // 0 = file does not exist yet
	}{
		{"new file", 0},
		{"existing world-readable file", 0644},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "sub", "secret.key")
			if tt.existing != 0 {
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("old content"), tt.existing); err != nil {
					t.Fatal(err)
				}
			}

			if err := WritePrivateFile(path, []byte("secret")); err != nil {
				t.Fatalf("WritePrivateFile() error = %v", err)
			}

			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if perm := info.Mode().Perm(); perm != 0600 {
				t.Errorf("permissions = %o, want 600", perm)
			}

			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != "secret" {
				t.Errorf("content = %q, want %q", data, "secret")
			}
		})
	}
}
