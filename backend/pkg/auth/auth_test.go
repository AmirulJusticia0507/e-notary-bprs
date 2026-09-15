package auth

import (
	"testing"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
)

func TestGenerateAndParseToken(t *testing.T) {
	user := &domain.User{ID: 42, Role: "admin"}

	token, err := GenerateToken(user, "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	claims, err := ParseToken(token, "test-secret")
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != user.ID {
		t.Errorf("UserID = %d, want %d", claims.UserID, user.ID)
	}
	if claims.Role != user.Role {
		t.Errorf("Role = %q, want %q", claims.Role, user.Role)
	}
}

func TestParseTokenInvalidSecret(t *testing.T) {
	user := &domain.User{ID: 1, Role: "legal_officer"}
	token, err := GenerateToken(user, "correct", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if _, err := ParseToken(token, "wrong"); err == nil {
		t.Error("ParseToken() with wrong secret should fail")
	}
}

func TestPasswordHashRoundtrip(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if !CheckPassword(hash, "secret123") {
		t.Error("CheckPassword() = false, want true")
	}
	if CheckPassword(hash, "wrong") {
		t.Error("CheckPassword() with wrong password = true, want false")
	}
}
