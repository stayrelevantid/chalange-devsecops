package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stayrelevantid/securebank-api/internal/middleware"
	"github.com/stayrelevantid/securebank-api/pkg/crypto"
)

var testJWTSecret = []byte("test-secret-key-for-testing-only")

func generateTestToken(t *testing.T) string {
	t.Helper()
	token, err := crypto.GenerateToken(testJWTSecret, "test-user", 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return token
}

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	healthCheck(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["status"] != "healthy" {
		t.Errorf("expected status healthy, got %s", resp["status"])
	}
}

func TestHealthCheckSecurityHeaders(t *testing.T) {
	handler := middleware.SecurityHeaders(middleware.LimitBodySize(1024, healthCheck))
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	checks := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Cache-Control":          "no-cache, no-store, must-revalidate, private",
	}
	for header, expected := range checks {
		got := w.Header().Get(header)
		if got != expected {
			t.Errorf("expected %s=%s, got %s", header, expected, got)
		}
	}
}

func TestGetBalanceFound(t *testing.T) {
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(1024,
			middleware.RequireAuth(testJWTSecret, getBalance),
		),
	)
	req := httptest.NewRequest("GET", "/balance?id=ACC001", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var acc Account
	if err := json.NewDecoder(w.Body).Decode(&acc); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if acc.ID != "ACC001" {
		t.Errorf("expected ACC001, got %s", acc.ID)
	}
}

func TestGetBalanceNotFound(t *testing.T) {
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(1024,
			middleware.RequireAuth(testJWTSecret, getBalance),
		),
	)
	req := httptest.NewRequest("GET", "/balance?id=INVALID", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetBalanceMissingID(t *testing.T) {
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(1024,
			middleware.RequireAuth(testJWTSecret, getBalance),
		),
	)
	req := httptest.NewRequest("GET", "/balance", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetBalanceNoAuth(t *testing.T) {
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(1024,
			middleware.RequireAuth(testJWTSecret, getBalance),
		),
	)
	req := httptest.NewRequest("GET", "/balance?id=ACC001", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetBalanceInvalidToken(t *testing.T) {
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(1024,
			middleware.RequireAuth(testJWTSecret, getBalance),
		),
	)
	req := httptest.NewRequest("GET", "/balance?id=ACC001", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestTransferSuccess(t *testing.T) {
	resetAccounts()
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(4096,
			middleware.RequireAuth(testJWTSecret, transfer),
		),
	)
	body, _ := json.Marshal(TransferReq{From: "ACC001", To: "ACC002", Amount: 500})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTransferInsufficientBalance(t *testing.T) {
	resetAccounts()
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(4096,
			middleware.RequireAuth(testJWTSecret, transfer),
		),
	)
	body, _ := json.Marshal(TransferReq{From: "ACC001", To: "ACC002", Amount: 999999})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTransferAccountNotFound(t *testing.T) {
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(4096,
			middleware.RequireAuth(testJWTSecret, transfer),
		),
	)
	body, _ := json.Marshal(TransferReq{From: "INVALID", To: "ACC002", Amount: 100})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestTransferNoAuth(t *testing.T) {
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(4096,
			middleware.RequireAuth(testJWTSecret, transfer),
		),
	)
	body, _ := json.Marshal(TransferReq{From: "ACC001", To: "ACC002", Amount: 100})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestTransferNegativeAmount(t *testing.T) {
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(4096,
			middleware.RequireAuth(testJWTSecret, transfer),
		),
	)
	body, _ := json.Marshal(TransferReq{From: "ACC001", To: "ACC002", Amount: -100})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative amount, got %d", w.Code)
	}
}

func TestTransferZeroAmount(t *testing.T) {
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(4096,
			middleware.RequireAuth(testJWTSecret, transfer),
		),
	)
	body, _ := json.Marshal(TransferReq{From: "ACC001", To: "ACC002", Amount: 0})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero amount, got %d", w.Code)
	}
}

func TestTransferEmptyAccountID(t *testing.T) {
	token := generateTestToken(t)
	handler := middleware.SecurityHeaders(
		middleware.LimitBodySize(4096,
			middleware.RequireAuth(testJWTSecret, transfer),
		),
	)
	body, _ := json.Marshal(TransferReq{From: "", To: "ACC002", Amount: 100})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty from account, got %d", w.Code)
	}
}

func resetAccounts() {
	mu.Lock()
	defer mu.Unlock()
	accounts["ACC001"] = &Account{ID: "ACC001", Name: "Alice", Balance: 10000}
	accounts["ACC002"] = &Account{ID: "ACC002", Name: "Bob", Balance: 5000}
}

func TestNotFoundHasSecurityHeaders(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", middleware.LimitBodySize(1024, healthCheck))
	handler := middleware.SecurityHeadersHandler(mux)

	req := httptest.NewRequest("GET", "/undefined-path", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	checks := map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"Cache-Control":           "no-cache, no-store, must-revalidate, private",
		"Content-Security-Policy": "default-src 'none'",
		"Referrer-Policy":         "strict-origin-when-cross-origin",
		"Permissions-Policy":      "camera=(), microphone=(), geolocation=()",
	}
	for header, expected := range checks {
		got := w.Header().Get(header)
		if got != expected {
			t.Errorf("expected %s=%s, got %s", header, expected, got)
		}
	}
}
