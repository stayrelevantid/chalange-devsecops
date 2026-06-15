package configs

import (
	"os"
	"testing"
)

func TestLoadDefaultValues(t *testing.T) {
	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("expected default Port 8080, got %s", cfg.Port)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected default DBHost localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPassword != "" {
		t.Errorf("expected default DBPassword empty, got %s", cfg.DBPassword)
	}
	if cfg.JWTSecret != "dev-secret-change-in-production" {
		t.Errorf("expected default JWTSecret dev-secret-change-in-production, got %s", cfg.JWTSecret)
	}
}

func TestLoadFromEnvVars(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PASSWORD", "mysecretpassword")
	os.Setenv("JWT_SECRET", "my-jwt-secret")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PASSWORD")
	defer os.Unsetenv("JWT_SECRET")

	cfg := Load()
	if cfg.Port != "9090" {
		t.Errorf("expected Port 9090, got %s", cfg.Port)
	}
	if cfg.DBHost != "db.example.com" {
		t.Errorf("expected DBHost db.example.com, got %s", cfg.DBHost)
	}
	if cfg.DBPassword != "mysecretpassword" {
		t.Errorf("expected DBPassword mysecretpassword, got %s", cfg.DBPassword)
	}
	if cfg.JWTSecret != "my-jwt-secret" {
		t.Errorf("expected JWTSecret my-jwt-secret, got %s", cfg.JWTSecret)
	}
}
