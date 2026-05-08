// --- COMMIT 2: Backend Base e Estruturas ---
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Estruturas de Dados
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Task struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil { log.Fatal(err) }

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", corsMiddleware(registerHandler))
	mux.HandleFunc("/api/v1/login", corsMiddleware(loginHandler))
	mux.HandleFunc("/api/v1/tasks", corsMiddleware(authMiddleware(tasksHandler)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Servidor Go rodando na porta 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Encerrando o servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao encerrar servidor: %v", err)
	}
	log.Println("Servidor encerrado com sucesso.")
}

// --- COMMIT 9: Handlers de Autenticação (Register/Login) ---
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { return }
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	_, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", user.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// JWT
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { return }
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	var storedHash string
	var userID int
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", user.Username).Scan(&userID, &storedHash)
	
	// Validação de senha
	if err != nil || bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password)) != nil {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// Geração do Token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// --- COMMIT 10: Handlers de CRUD de Tarefas ---
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)

	if r.Method == http.MethodGet {
		rows, _ := db.Query("SELECT id, title, completed FROM tasks WHERE user_id = ?", userID)
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var t Task
			rows.Scan(&t.ID, &t.Title, &t.Completed)
			tasks = append(tasks, t)
		}
		json.NewEncoder(w).Encode(tasks)
	} else if r.Method == http.MethodPost {
		var t Task
		json.NewDecoder(r.Body).Decode(&t)
		db.Exec("INSERT INTO tasks (user_id, title) VALUES (?, ?)", userID, t.Title)
		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodPut {
		var t Task
		json.NewDecoder(r.Body).Decode(&t)
		_, err := db.Exec("UPDATE tasks SET title = ?, completed = ? WHERE id = ? AND user_id = ?", t.Title, t.Completed, t.ID, userID)
		if err != nil {
			http.Error(w, "Erro ao atualizar tarefa", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodDelete {
		id := r.URL.Query().Get("id")
		_, err := db.Exec("DELETE FROM tasks WHERE id = ? AND user_id = ?", id, userID)
		if err != nil {
			http.Error(w, "Erro ao deletar tarefa", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// --- COMMIT 11: Middlewares (CORS e Segurança) ---
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) > 7 { tokenString = tokenString[7:] } // Remove "Bearer "

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", claims["user_id"])
			next(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
		}
	}
}