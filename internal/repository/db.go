package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type DB struct {
	Conn *pgx.Conn
}

func NewDB(ctx context.Context) (*DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{Conn: conn}, nil
}

func (db *DB) Close(ctx context.Context) error {
	return db.Conn.Close(ctx)
}
