package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"

	"github.com/stayrelevantid/securebank-api/configs"
	"github.com/stayrelevantid/securebank-api/internal/middleware"
)

type Account struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type TransferReq struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	accounts = map[string]*Account{
		"ACC001": {ID: "ACC001", Name: "Alice", Balance: 10000},
		"ACC002": {ID: "ACC002", Name: "Bob", Balance: 5000},
	}
	mu sync.RWMutex
)

const maxTransferAmount = 1_000_000_000.00

func getBalance(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing account id"})
		return
	}
	mu.RLock()
	acc, ok := accounts[id]
	mu.RUnlock()
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "account not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func transfer(w http.ResponseWriter, r *http.Request) {
	var req TransferReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	req.From = strings.TrimSpace(req.From)
	req.To = strings.TrimSpace(req.To)

	if req.From == "" || req.To == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "from and to account ids are required"})
		return
	}
	if req.Amount <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "transfer amount must be positive"})
		return
	}
	if req.Amount > maxTransferAmount {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "transfer amount exceeds maximum limit"})
		return
	}
	if math.IsInf(req.Amount, 0) || math.IsNaN(req.Amount) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid transfer amount"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	from, ok1 := accounts[req.From]
	to, ok2 := accounts[req.To]
	if !ok1 || !ok2 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "account not found"})
		return
	}
	if from.Balance < req.Amount {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "insufficient balance"})
		return
	}
	from.Balance -= req.Amount
	to.Balance += req.Amount

	log.Printf("TRANSFER from=%s to=%s amount=%.2f", req.From, req.To, req.Amount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
	cfg := configs.Load()
	jwtSecret := []byte(cfg.JWTSecret)

	authBalance := middleware.SecurityHeaders(
		middleware.LimitBodySize(1024,
			middleware.RequireAuth(jwtSecret, getBalance),
		),
	)

	authTransfer := middleware.SecurityHeaders(
		middleware.LimitBodySize(4096,
			middleware.RequireAuth(jwtSecret, transfer),
		),
	)

	publicHealth := middleware.SecurityHeaders(
		middleware.LimitBodySize(1024, healthCheck),
	)

	http.HandleFunc("/health", publicHealth)
	http.HandleFunc("/balance", authBalance)
	http.HandleFunc("/transfer", authTransfer)
	log.Printf("SecureBank API running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
