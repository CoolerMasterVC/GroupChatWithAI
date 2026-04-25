package main

import (
	"log"
	"net/http"
	"time"

	"transport/internal/handlers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "transport/docs" // для swagger
)

// @title Transport Layer API (Stub)
// @version 1.0
// @description Заглушка транспортного уровня для демонстрации Swagger.
// @host localhost:8080
// @BasePath /
func main() {
	r := mux.NewRouter()

	// Эндпоинт для приёма сегментов
	r.HandleFunc("/segment", handlers.ReceiveSegment).Methods("POST")

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Простой middleware для CORS (если нужно)
	r.Use(corsMiddleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server starting on :8080")
	log.Fatal(srv.ListenAndServe())
}

// corsMiddleware добавляет заголовки CORS
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
