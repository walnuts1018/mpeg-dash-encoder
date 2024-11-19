package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func ReaderHash(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", fmt.Errorf("failed to copy: %w", err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
