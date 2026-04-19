package encodingutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func Hash(algo string, data []byte) ([]byte, error) {
	switch algo {
	case "md5":
		sum := md5.Sum(data)
		return sum[:], nil

	case "sha1":
		sum := sha1.Sum(data)
		return sum[:], nil

	case "sha256":
		sum := sha256.Sum256(data)
		return sum[:], nil

	case "sha512":
		sum := sha512.Sum512(data)
		return sum[:], nil

	default:
		return nil, fmt.Errorf("unsupported algorithm")
	}
}
