package repository

import (
	"errors"
	"fmt"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	WalletNotFound       = errors.New("wallet not found")
)

func applyOperation(balance int64, operationType string, amount int64) (int64, error) {
	switch operationType {
	case "DEPOSIT":
		balance += amount
	case "WITHDRAW":
		if balance < amount {
			return 0, ErrInsufficientFunds
		}
		balance -= amount
	default:
		return 0, fmt.Errorf("invalid operation type: %s", operationType)
	}
	return balance, nil
}
