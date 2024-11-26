// main.go
package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
	"wallet-service/internal/api"
	"wallet-service/internal/repository"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load("/app/config/config.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Устанавливаем тайм-аут
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Создаем подключение к базе данных
	db, err := repository.NewDB(ctx)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close(ctx)

	// Создаем сервер
	server := api.NewServer(db)

	// Настройка маршрутов
	r := mux.NewRouter()

	// Регистрация маршрутов
	r.HandleFunc("/api/v1/wallet", server.HandleTransaction).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/wallets/{id}", server.HandleGetBalance).Methods(http.MethodGet)

	// Получаем порт из конфигурации
	port := os.Getenv("APP_PORT")
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
