package formatutil

import "testing"

func TestFormatLogicalSize(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{999, "999 B"},
		{1000, "1.0 KB"},
		{1500, "1.5 KB"},
		{999999, "1000.0 KB"},
		{1000000, "1.0 MB"},
		{1500000, "1.5 MB"},
		{999999999, "1000.0 MB"},
		{1000000000, "1.0 GB"},
		{1500000000, "1.5 GB"},
		{10000000000, "10.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatLogicalSize(tt.input)
			if result != tt.expected {
				t.Errorf("FormatLogicalSize(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
