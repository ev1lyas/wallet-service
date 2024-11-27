package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"net/http"
	"wallet-service/internal/repository"
)

// HandleTransaction обрабатывает операции с кошельком
func (s *Server) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Проверка на корректность UUID
	_, err := uuid.Parse(req.WalletId)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "Invalid walletId format, must be a valid UUID")
		return
	}

	// Проверка, что сумма больше нуля
	if req.Amount <= 0 {
		ResponseError(w, http.StatusBadRequest, "Amount must be greater than 0")
		return
	}

	// Проверка на допустимые операции
	if req.OperationType != "DEPOSIT" && req.OperationType != "WITHDRAW" {
		ResponseError(w, http.StatusBadRequest, "Invalid operation type. Supported types are DEPOSIT and WITHDRAW")
		return
	}

	// Обновление баланса
	ctx := r.Context()
	if err := s.DB.UpdateBalance(ctx, req.WalletId, req.OperationType, req.Amount); err != nil {
		if err == repository.ErrInsufficientFunds {
			ResponseError(w, http.StatusConflict, "Insufficient funds")
			return
		}
		if err == pgx.ErrNoRows {
			ResponseError(w, http.StatusNotFound, "Wallet not found")
			return
		}
		ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleGetBalance обрабатывает запросы на получение баланса кошелька
func (s *Server) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID, ok := vars["id"]
	if !ok || walletID == "" {
		ResponseError(w, http.StatusBadRequest, "Wallet ID is required")
		return
	}

	ctx := r.Context()
	balance, err := s.DB.GetBalance(ctx, walletID)
	if err != nil {
		if err.Error() == "wallet not found" || err == pgx.ErrNoRows {
			ResponseError(w, http.StatusNotFound, "Wallet not found")
			return
		}
		ResponseError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := ResponseBody{
		WalletID: walletID,
		Balance:  balance,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ResponseError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}
