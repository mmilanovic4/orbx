package encodingutil

import (
	"bytes"
	"testing"
)

func TestBase64RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"empty", []byte{}},
		{"hello", []byte("Hello!")},
		{"binary", []byte{0x00, 0xFF, 0xAB, 0x10}},
		{"unicode", []byte("こんにちは")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeBase64(tt.input)
			decoded, err := DecodeBase64(encoded)
			if err != nil {
				t.Fatalf("DecodeBase64 failed: %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("got %q, want %q", decoded, tt.input)
			}
		})
	}
}

func TestDecodeBase64Invalid(t *testing.T) {
	_, err := DecodeBase64("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestDecodeBase64URL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"url alphabet", "PD4_fn5-", "<>?~~~"},
		{"standard alphabet", "PD4/fn5+", "<>?~~~"},
		{"without padding", "SGk", "Hi"},
		{"with padding", "SGk=", "Hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := DecodeBase64URL(tt.input)
			if err != nil {
				t.Fatalf("DecodeBase64URL(%q) error = %v", tt.input, err)
			}
			if string(decoded) != tt.expected {
				t.Errorf("got %q, want %q", decoded, tt.expected)
			}
		})
	}
}

func TestDecodeBase64URLInvalid(t *testing.T) {
	for _, input := range []string{"a", "not valid!"} {
		if _, err := DecodeBase64URL(input); err == nil {
			t.Errorf("DecodeBase64URL(%q) expected error, got nil", input)
		}
	}
}

func TestHexRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"empty", []byte{}},
		{"hello", []byte("hello")},
		{"binary", []byte{0x00, 0xFF, 0xAB}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeHex(tt.input)
			decoded, err := DecodeHex(encoded)
			if err != nil {
				t.Fatalf("DecodeHex failed: %v", err)
			}
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("got %q, want %q", decoded, tt.input)
			}
		})
	}
}

func TestDecodeHexInvalid(t *testing.T) {
	_, err := DecodeHex("zzzz")
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestHash(t *testing.T) {
	tests := []struct {
		algo     string
		input    []byte
		expected string
	}{
		{"md5", []byte("hello"), "5d41402abc4b2a76b9719d911017c592"},
		{"sha1", []byte("hello"), "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{"sha256", []byte("hello"), "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"sha512", []byte("hello"), "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"},
	}

	for _, tt := range tests {
		t.Run(tt.algo, func(t *testing.T) {
			result, err := Hash(tt.algo, tt.input)
			if err != nil {
				t.Fatalf("Hash failed: %v", err)
			}
			if EncodeHex(result) != tt.expected {
				t.Errorf("got %s, want %s", EncodeHex(result), tt.expected)
			}
		})
	}
}

func TestHashUnsupportedAlgo(t *testing.T) {
	_, err := Hash("md6", []byte("hello"))
	if err == nil {
		t.Error("expected error for unsupported algorithm")
	}
}

func TestHTMLRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		encoded string
	}{
		{"ampersand", "a & b", "a &amp; b"},
		{"tags", "<div>", "&lt;div&gt;"},
		{"quotes", `"hello"`, "&#34;hello&#34;"},
		{"plain", "hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeHTML(tt.input)
			if encoded != tt.encoded {
				t.Errorf("EncodeHTML got %q, want %q", encoded, tt.encoded)
			}
			decoded := DecodeHTML(encoded)
			if decoded != tt.input {
				t.Errorf("DecodeHTML got %q, want %q", decoded, tt.input)
			}
		})
	}
}
