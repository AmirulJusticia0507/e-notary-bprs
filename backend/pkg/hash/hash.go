package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func FromReader(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func FromBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func FromString(value string) string {
	return FromBytes([]byte(value))
}
