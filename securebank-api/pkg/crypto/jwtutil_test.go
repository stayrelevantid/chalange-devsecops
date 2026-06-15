package crypto

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := []byte("my-secret-key")
	token, err := GenerateToken(secret, "user123", 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}

	userID, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}
	if userID != "user123" {
		t.Errorf("expected userID user123, got %s", userID)
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	secret1 := []byte("secret-one")
	secret2 := []byte("secret-two")
	token, err := GenerateToken(secret1, "user123", 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	_, err = ParseToken(secret2, token)
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}

func TestParseTokenExpired(t *testing.T) {
	secret := []byte("my-secret-key")
	token, err := GenerateToken(secret, "user123", -1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	_, err = ParseToken(secret, token)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestParseTokenInvalidString(t *testing.T) {
	secret := []byte("my-secret-key")
	_, err := ParseToken(secret, "not-a-valid-token")
	if err == nil {
		t.Error("expected error for invalid token string, got nil")
	}
}

func TestHashPasswordAndCheck(t *testing.T) {
	hash, err := HashPassword("mysecretpassword")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if !CheckPassword("mysecretpassword", hash) {
		t.Error("CheckPassword should return true for correct password")
	}
	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword should return false for incorrect password")
	}
}
