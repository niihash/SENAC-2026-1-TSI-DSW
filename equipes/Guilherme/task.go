package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *sql.DB // Assumindo que a conexão já foi inicializada no main.go

// GET: Listar todas as tarefas
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, completed FROM tasks")
	if err != nil {
		http.Error(w, "Erro ao buscar tarefas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// POST: Criar uma nova tarefa
func createTasksHandler(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO tasks(title, completed) VALUES(?, ?)")
	if err != nil {
		http.Error(w, "Erro no servidor", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Title, false)
	if err != nil {
		http.Error(w, "Erro ao inserir no banco", http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	t.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

// PUT: Atualizar status da tarefa (Concluir)
func updateTasksHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	id, _ := strconv.Atoi(idStr)

	var t Task
	json.NewDecoder(r.Body).Decode(&t)

	stmt, err := db.Prepare("UPDATE tasks SET completed = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Erro no servidor", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Completed, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
