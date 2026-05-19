package main

import (
	"log"
	"net/http"
	"time"

	"transport/internal/assembler"
	"transport/internal/handlers"
	"transport/internal/kafka"
	"transport/internal/storage"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "transport/docs"
)

// @title Transport Layer API
// @version 1.0
// @description API транспортного уровня с Kafka
// @host localhost:8080
// @BasePath /
func main() {
	store := storage.NewStorage()

	// Запуск Kafka consumer (в фоне)
	kafka.StartConsumer(store)
	defer kafka.StopConsumer()

	// Запуск сборщика
	assembler.StartAssembler(store)

	r := mux.NewRouter()

	// /send не требует store
	r.HandleFunc("/send", handlers.HandleSend).Methods("POST", "OPTIONS")

	// /transfer требует store – используем замыкание
	r.HandleFunc("/transfer", handlers.HandleTransfer(store)).Methods("POST", "OPTIONS")

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.Use(corsMiddleware)

	// Слушаем все интерфейсы (0.0.0.0:8080) – это обеспечит доступ по ZeroTier
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server starting on :8080 (all interfaces)")
	log.Fatal(srv.ListenAndServe())
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
