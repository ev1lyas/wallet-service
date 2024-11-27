package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"wallet-service/internal/api"
	"wallet-service/internal/repository" // Импортируем репозиторий
)

func setupTestDB(t *testing.T) *repository.DB {
	t.Helper()
	// Создаём подключение к тестовой базе данных
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Создаем обертку DB
	db := &repository.DB{Conn: conn}

	// Применяем миграции
	err = applyTestMigrations(context.Background(), conn)
	if err != nil {
		conn.Close(context.Background()) // Закрываем соединение в случае ошибки
		t.Fatalf("failed to apply migrations: %v", err)
	}

	return db
}

// applyTestMigrations применяет миграции к тестовой базе данных
func applyTestMigrations(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS wallets (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			balance BIGINT NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create test schema: %v", err)
	}
	return nil
}

// TestHandleTransaction проверяет обработку транзакций (депозит/снятие)
func TestHandleTransaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close(context.Background())

	// Создаём сервер с обработчиками
	server := api.NewServer(db) // Теперь передаем db (объект repository.DB)

	// Подготовка данных для POST запроса
	payload := map[string]interface{}{
		"wallet_id":      "test-wallet-id",
		"operation_type": "deposit",
		"amount":         500,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Router.ServeHTTP(rec, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
}

// TestHandleGetBalance проверяет получение баланса кошелька
func TestHandleGetBalance(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close(context.Background())

	// Вставляем тестовые данные
	_, err := db.Conn.Exec(context.Background(), `
		INSERT INTO wallets (id, balance) VALUES ('test-wallet-id', 1000)
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Создаём сервер с обработчиками
	server := api.NewServer(db) // Теперь передаем db (объект repository.DB)

	// Создаём GET запрос
	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/test-wallet-id", nil)
	rec := httptest.NewRecorder()
	server.Router.ServeHTTP(rec, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, response["balance"])
}
