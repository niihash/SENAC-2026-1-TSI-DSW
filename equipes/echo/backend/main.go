package main

import (
	"backend/controller"
	"backend/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	database.Connect()

	defer database.DB.Close()

	r := mux.NewRouter()

	// tasks
	r.HandleFunc("/tasks", controller.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", controller.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", controller.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", controller.DeleteTask).Methods("DELETE")

	// rota auth(users)
	r.HandleFunc("/register", controller.Register).Methods("POST")
	r.HandleFunc("/login", controller.Login).Methods("POST")
	r.HandleFunc("/logout", controller.Logout).Methods("POST")

	// Melhor maneira: define o arquivo principal na raiz para nao listar diretorio
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/login.html")
	})

	// Esta rota deve vir por ultimo. Ela serve os arquivos da pasta frontend
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("frontend/")))

	log.Println("Servidor rodando na porta 8080")

	// "0.0.0.0" para o container aceitar conexões de fora (navegador)
	http.ListenAndServe("0.0.0.0:8080", enableCORS(r))
}
