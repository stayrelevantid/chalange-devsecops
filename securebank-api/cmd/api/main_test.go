package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestGetBalanceFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/balance?id=ACC001", nil)
	w := httptest.NewRecorder()
	getBalance(w, req)
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
	req := httptest.NewRequest("GET", "/balance?id=INVALID", nil)
	w := httptest.NewRecorder()
	getBalance(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestTransferSuccess(t *testing.T) {
	body, _ := json.Marshal(TransferReq{From: "ACC001", To: "ACC002", Amount: 500})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	transfer(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTransferInsufficientBalance(t *testing.T) {
	body, _ := json.Marshal(TransferReq{From: "ACC001", To: "ACC002", Amount: 999999})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	transfer(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTransferAccountNotFound(t *testing.T) {
	body, _ := json.Marshal(TransferReq{From: "INVALID", To: "ACC002", Amount: 100})
	req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	transfer(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
