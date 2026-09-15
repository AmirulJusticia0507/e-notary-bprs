package pdf

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(_ context.Context, _ string, _ any) ([]byte, error) {
	return nil, errors.New("PDF generator is not configured")
}

func (g *Generator) Hash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
