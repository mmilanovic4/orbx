package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"orbx/internal/encodingutil"
	"orbx/internal/sysutil"
	"strings"
)

const (
	KeySize128 = 16
	KeySize192 = 24
	KeySize256 = 32
)

func validateKeySize(key []byte) error {
	switch len(key) {
	case KeySize128, KeySize192, KeySize256:
		return nil
	default:
		return fmt.Errorf("invalid key size %d: must be 16, 24, or 32 bytes", len(key))
	}
}

func newGCM(key []byte) (cipher.AEAD, error) {
	if err := validateKeySize(key); err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	return gcm, nil
}

func Encrypt(plaintext, key []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(data, key []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize+gcm.Overhead() {
		return nil, fmt.Errorf("invalid ciphertext: data too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

func ReadKey(keyFile string) ([]byte, error) {
	if keyFile == "" {
		return nil, fmt.Errorf("key file path not provided")
	}

	keyEncoded, err := sysutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	trimmed := strings.TrimSpace(string(keyEncoded))

	keyDecoded, err := encodingutil.DecodeBase64(trimmed)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	if err := validateKeySize(keyDecoded); err != nil {
		return nil, fmt.Errorf("invalid key in file: %w", err)
	}

	return keyDecoded, nil
}
