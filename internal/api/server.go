package api

import (
	"github.com/gorilla/mux"
	"wallet-service/internal/repository"
)

// Server описывает сервер с роутером и базой данных
type Server struct {
	Router *mux.Router
	DB     *repository.DB
}

// NewServer создает новый сервер
func NewServer(db *repository.DB) *Server {
	server := &Server{DB: db}
	server.Router = mux.NewRouter()

	// Регистрируем маршруты
	server.Router.HandleFunc("/api/v1/wallet", server.HandleTransaction).Methods("POST")
	server.Router.HandleFunc("/api/v1/wallets/{id}", server.HandleGetBalance).Methods("GET")
	return server
}
