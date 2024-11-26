package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (db *DB) UpdateBalance(ctx context.Context, walletId, operationType string, amount int64) error {
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	balance, err := db.getWalletBalance(ctx, tx, walletId)
	if err != nil {
		if err == pgx.ErrNoRows {
			balance, err = db.createWallet(ctx, tx, walletId)
			if err != nil {
				return fmt.Errorf("failed to create wallet: %w", err)
			}
		} else {
			return fmt.Errorf("failed to retrieve wallet balance: %w", err)
		}
	}

	balance, err = applyOperation(balance, operationType, amount)
	if err != nil {
		return fmt.Errorf("failed to apply operation: %w", err)
	}

	if err := db.updateWalletBalance(ctx, tx, walletId, balance); err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	return nil
}

func (db *DB) GetBalance(ctx context.Context, walletId string) (int64, error) {
	var balance int64
	err := db.Conn.QueryRow(ctx, `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`, walletId).Scan(&balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, WalletNotFound
		}
		return 0, fmt.Errorf("failed to retrieve balance: %w", err)
	}
	return balance, nil
}
