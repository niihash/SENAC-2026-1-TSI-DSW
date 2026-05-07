package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
)

type App struct {
	db *sql.DB
}

func main() {
	db := newDB()
	defer db.Close()

	app := &App{db: db}

	mux := http.NewServeMux()

	// rotas de tasks
	mux.HandleFunc("/api/v1/tasks", app.tasksCollectionHandler)
	mux.HandleFunc("/api/v1/tasks/", app.taskItemHandler)

	port := getEnv("PORT", "8080")
	log.Printf("servidor rodando na porta %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}

// /api/v1/tasks  →  GET (listar) | POST (criar)
func (app *App) tasksCollectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.listTasks(w, r)
	case http.MethodPost:
		app.createTask(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
	}
}

// /api/v1/tasks/{id}  →  GET | PUT | DELETE
func (app *App) taskItemHandler(w http.ResponseWriter, r *http.Request) {
	// garante que não cai aqui com path exatamente "/api/v1/tasks/"
	idStr := extractIDFromPath(r.URL.Path)
	if idStr == "" || strings.TrimSpace(idStr) == "" {
		writeError(w, http.StatusBadRequest, "id não informado")
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.getTask(w, r)
	case http.MethodPut:
		app.updateTask(w, r)
	case http.MethodDelete:
		app.deleteTask(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
	}
}