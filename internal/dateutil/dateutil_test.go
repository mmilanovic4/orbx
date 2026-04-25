package dateutil

import (
	"testing"
	"time"
)

func TestGetLocation(t *testing.T) {
	tests := []struct {
		name    string
		tz      string
		wantErr bool
	}{
		{"empty returns local", "", false},
		{"valid timezone", "Europe/Belgrade", false},
		{"valid UTC", "UTC", false},
		{"invalid timezone", "Invalid/Zone", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := GetLocation(tt.tz)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocation(%q) error = %v, wantErr = %v", tt.tz, err, tt.wantErr)
			}
			if !tt.wantErr && loc == nil {
				t.Error("expected non-nil location")
			}
			if tt.tz == "" && loc != time.Local {
				t.Error("expected time.Local for empty timezone")
			}
		})
	}
}

func TestGetLayout(t *testing.T) {
	tests := []struct {
		name     string
		ms       bool
		expected string
	}{
		{"without ms", false, "2006-01-02 15:04:05"},
		{"with ms", true, "2006-01-02 15:04:05.000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLayout(tt.ms)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}
