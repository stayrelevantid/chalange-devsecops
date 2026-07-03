package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stayrelevantid/securebank-api/pkg/crypto"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func TestSecurityHeaders(t *testing.T) {
	handler := SecurityHeaders(http.HandlerFunc(okHandler))
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	checks := map[string]string{
		"X-Content-Type-Options":       "nosniff",
		"X-Frame-Options":              "DENY",
		"Cache-Control":                "no-cache, no-store, must-revalidate, private",
		"Content-Security-Policy":      "default-src 'none'",
		"X-XSS-Protection":             "0",
		"Referrer-Policy":              "strict-origin-when-cross-origin",
		"Permissions-Policy":           "camera=(), microphone=(), geolocation=()",
		"Cross-Origin-Resource-Policy": "same-origin",
	}
	for header, expected := range checks {
		got := w.Header().Get(header)
		if got != expected {
			t.Errorf("SecurityHeaders: expected %s=%s, got %s", header, expected, got)
		}
	}
}

func TestLimitBodySizeAllowed(t *testing.T) {
	handler := LimitBodySize(1024, http.HandlerFunc(okHandler))
	body := strings.NewReader(`{"data":"small"}`)
	req := httptest.NewRequest("POST", "/test", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestLimitBodySizeRejected(t *testing.T) {
	handler := LimitBodySize(10, http.HandlerFunc(okHandler))
	largeBody := strings.NewReader(strings.Repeat("a", 100))
	req := httptest.NewRequest("POST", "/test", largeBody)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413, got %d", w.Code)
	}
}

func TestRequireAuthNoHeader(t *testing.T) {
	secret := []byte("test-secret")
	handler := RequireAuth(secret, http.HandlerFunc(okHandler))
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without auth header, got %d", w.Code)
	}
}

func TestRequireAuthInvalidFormat(t *testing.T) {
	secret := []byte("test-secret")
	handler := RequireAuth(secret, http.HandlerFunc(okHandler))
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic abc123")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for non-Bearer format, got %d", w.Code)
	}
}

func TestRequireAuthInvalidToken(t *testing.T) {
	secret := []byte("test-secret")
	handler := RequireAuth(secret, http.HandlerFunc(okHandler))
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for invalid token, got %d", w.Code)
	}
}

func TestRequireAuthValidToken(t *testing.T) {
	secret := []byte("test-secret")
	token, err := crypto.GenerateToken(secret, "user1", 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	handler := RequireAuth(secret, http.HandlerFunc(okHandler))
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for valid token, got %d", w.Code)
	}
}

func TestRequireAuthWrongSecret(t *testing.T) {
	secret1 := []byte("secret-one")
	secret2 := []byte("secret-two")
	token, err := crypto.GenerateToken(secret1, "user1", 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	handler := RequireAuth(secret2, http.HandlerFunc(okHandler))
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong secret, got %d", w.Code)
	}
}
