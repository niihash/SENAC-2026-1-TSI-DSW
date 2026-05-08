package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const serverAddr = ":8080"

func main() {
	var err error

	_ = godotenv.Load()

	db, err = sql.Open("mysql", databaseDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	frontendDir := filepath.Join("..", "frontend")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/tasks", tasksCollectionHandler)
	mux.HandleFunc("/api/v1/tasks/", tasksItemHandler)
	mux.Handle("/", http.FileServer(http.Dir(frontendDir)))

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      securityHeaders(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Servidor rodando em http://localhost%s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func databaseDSN() string {
	user := envOrDefault("DB_USER", "root")
	password := os.Getenv("DB_PASSWORD")
	host := envOrDefault("DB_HOST", "127.0.0.1")
	port := envOrDefault("DB_PORT", "3306")
	name := envOrDefault("DB_NAME", "todo_list")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, name)
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func tasksCollectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	case http.MethodGet:
		getTasksHandler(w, r)
	case http.MethodPost:
		createTasksHandler(w, r)
	default:
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
	}
}

func tasksItemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPut:
		updateTasksHandler(w, r)
	case http.MethodDelete:
		deleteTasksHandler(w, r)
	default:
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	allowedOrigin := envOrDefault("CORS_ALLOWED_ORIGIN", "http://localhost:8080")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Header.Get("Origin") == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		next.ServeHTTP(w, r)
	})
}
