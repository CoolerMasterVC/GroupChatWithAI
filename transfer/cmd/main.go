package main

import (
	"log"
	"net/http"
	"time"

	"transport/internal/handlers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "transport/docs" // Сгенерированная документация Swagger
)

// @title Transport Layer API
// @version 1.0
// @description API транспортного уровня для обмена сообщениями между прикладным и агентным уровнями.
// @host localhost:8080
// @BasePath /
func main() {
	r := mux.NewRouter()

	// Эндпоинты
	r.HandleFunc("/send", handlers.HandleSend).Methods("POST")
	r.HandleFunc("/transfer", handlers.HandleTransfer).Methods("POST")

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// CORS (опционально)
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
