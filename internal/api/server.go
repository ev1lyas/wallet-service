package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"net/http"
	"wallet-service/internal/repository"
)

type Server struct {
	Router *mux.Router
	DB     *repository.DB
}

func NewServer(db *repository.DB) *Server {
	return &Server{DB: db}
}

// HandleTransaction обрабатывает операции с кошельком
func (s *Server) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		WalletId      string `json:"walletId"`
		OperationType string `json:"operationType"`
		Amount        int64  `json:"amount"`
	}

	// Декодируем тело запроса
	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка на корректность UUID
	_, err := uuid.Parse(req.WalletId)
	if err != nil {
		http.Error(w, "Invalid walletId format, must be a valid UUID", http.StatusBadRequest)
		return
	}

	// Проверка, что сумма больше нуля
	if req.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	// Обновление баланса
	ctx := r.Context()
	if err := s.DB.UpdateBalance(ctx, req.WalletId, req.OperationType, req.Amount); err != nil {
		if err == repository.ErrInsufficientFunds {
			http.Error(w, "Insufficient funds", http.StatusConflict)
			return
		}
		if err == pgx.ErrNoRows {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleGetBalance обрабатывает запросы на получение баланса кошелька
func (s *Server) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID кошелька из параметров маршрута
	vars := mux.Vars(r)
	walletID, ok := vars["id"]
	if !ok || walletID == "" {
		http.Error(w, "Wallet ID is required", http.StatusBadRequest)
		return
	}

	// Получаем баланс из базы данных
	ctx := r.Context()
	balance, err := s.DB.GetBalance(ctx, walletID)
	if err != nil {
		if err.Error() == "wallet not found" {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}
		if err == pgx.ErrNoRows {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Формируем ответ в формате JSON
	response := struct {
		WalletID string `json:"walletId"`
		Balance  int64  `json:"balance"`
	}{
		WalletID: walletID,
		Balance:  balance,
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
