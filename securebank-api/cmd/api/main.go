package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
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

var (
	accounts = map[string]*Account{
		"ACC001": {ID: "ACC001", Name: "Alice", Balance: 10000},
		"ACC002": {ID: "ACC002", Name: "Bob", Balance: 5000},
	}
	mu sync.RWMutex
)

func getBalance(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	mu.RLock()
	acc, ok := accounts[id]
	mu.RUnlock()
	if !ok {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func transfer(w http.ResponseWriter, r *http.Request) {
	var req TransferReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	from, ok1 := accounts[req.From]
	to, ok2 := accounts[req.To]
	if !ok1 || !ok2 {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}
	if from.Balance < req.Amount {
		http.Error(w, "insufficient balance", http.StatusBadRequest)
		return
	}
	from.Balance -= req.Amount
	to.Balance += req.Amount
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/balance", getBalance)
	http.HandleFunc("/transfer", transfer)
	log.Println("SecureBank API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
