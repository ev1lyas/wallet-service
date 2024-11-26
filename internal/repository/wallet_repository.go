package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func (db *DB) getWalletBalance(ctx context.Context, tx pgx.Tx, walletId string) (int64, error) {
	var balance int64
	err := tx.QueryRow(ctx, `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`, walletId).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (db *DB) createWallet(ctx context.Context, tx pgx.Tx, walletId string) (int64, error) {
	_, err := tx.Exec(ctx, `INSERT INTO wallets (id, balance) VALUES ($1, $2)`, walletId, 0)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (db *DB) updateWalletBalance(ctx context.Context, tx pgx.Tx, walletId string, balance int64) error {
	_, err := tx.Exec(ctx, `UPDATE wallets SET balance = $1 WHERE id = $2`, balance, walletId)
	return err
}
