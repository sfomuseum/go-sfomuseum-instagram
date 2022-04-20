package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
)

// FileHash returns the SHA-256 hash for the contents of 'r'.
func FileHash(r io.Reader) (string, error) {

	h := sha256.New()
	_, err := io.Copy(h, r)

	if err != nil {
		return "", fmt.Errorf("Failed to copy data, %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
