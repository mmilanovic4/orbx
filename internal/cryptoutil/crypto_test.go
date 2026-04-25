package cryptoutil

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func randomKey(size int) []byte {
	key := make([]byte, size)
	rand.Read(key)
	return key
}

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		input   []byte
	}{
		{"aes128", KeySize128, []byte("Hello from the other side!")},
		{"aes192", KeySize192, []byte("Hello from the other side!")},
		{"aes256", KeySize256, []byte("Hello from the other side!")},
		{"empty", KeySize256, []byte{}},
		{"binary", KeySize256, []byte{0x00, 0xFF, 0x10, 0xAB}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := randomKey(tt.keySize)

			ciphertext, err := Encrypt(tt.input, key)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			plaintext, err := Decrypt(ciphertext, key)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if !bytes.Equal(plaintext, tt.input) {
				t.Errorf("got %q, want %q", plaintext, tt.input)
			}
		})
	}
}

func TestEncryptIsNonDeterministic(t *testing.T) {
	key := randomKey(KeySize256)
	input := []byte("same input")

	a, _ := Encrypt(input, key)
	b, _ := Encrypt(input, key)

	if bytes.Equal(a, b) {
		t.Error("expected different ciphertexts for same input (nonce should differ)")
	}
}

func TestDecryptInvalidKey(t *testing.T) {
	key := randomKey(KeySize256)
	input := []byte("Hello!")

	ciphertext, _ := Encrypt(input, key)

	wrongKey := randomKey(KeySize256)
	_, err := Decrypt(ciphertext, wrongKey)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

func TestDecryptTamperedData(t *testing.T) {
	key := randomKey(KeySize256)
	ciphertext, _ := Encrypt([]byte("Hello!"), key)

	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err := Decrypt(ciphertext, key)
	if err == nil {
		t.Error("expected error when decrypting tampered ciphertext")
	}
}

func TestDecryptTooShort(t *testing.T) {
	key := randomKey(KeySize256)
	_, err := Decrypt([]byte("short"), key)
	if err == nil {
		t.Error("expected error for data that is too short")
	}
}

func TestValidateKeySize(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{"16 bytes", KeySize128, false},
		{"24 bytes", KeySize192, false},
		{"32 bytes", KeySize256, false},
		{"invalid 0", 0, true},
		{"invalid 10", 10, true},
		{"invalid 64", 64, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKeySize(make([]byte, tt.size))
			if (err != nil) != tt.wantErr {
				t.Errorf("validateKeySize(%d) error = %v, wantErr = %v", tt.size, err, tt.wantErr)
			}
		})
	}
}

func TestShannonEntropy(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantMin float64
		wantMax float64
	}{
		{"empty", []byte{}, 0, 0},
		{"single byte", []byte{0x00}, 0, 0},
		{"all same", []byte{0xAA, 0xAA, 0xAA}, 0, 0},
		{"two values", []byte{0x00, 0xFF, 0x00, 0xFF}, 0.9, 1.1},
		{"high entropy", func() []byte {
			b := make([]byte, 10000)
			rand.Read(b)
			return b
		}(), 7.9, 8.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShannonEntropy(tt.input)
			if result < tt.wantMin || result > tt.wantMax {
				t.Errorf("ShannonEntropy() = %.4f, want between %.4f and %.4f", result, tt.wantMin, tt.wantMax)
			}
		})
	}
}
