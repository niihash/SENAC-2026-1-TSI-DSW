package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *sql.DB // Inicializada no main.go

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, completed FROM tasks ORDER BY id")
	if err != nil {
		http.Error(w, "Erro ao buscar tarefas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
			http.Error(w, "Erro ao ler tarefas", http.StatusInternalServerError)
			return
		}

		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTasksHandler(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	t.Title = strings.TrimSpace(t.Title)
	t.Completed = false
	if t.Title == "" {
		http.Error(w, "Titulo obrigatorio", http.StatusBadRequest)
		return
	}
	if len([]rune(t.Title)) > 255 {
		http.Error(w, "Titulo deve ter no maximo 255 caracteres", http.StatusBadRequest)
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
	t.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func updateTasksHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE tasks SET completed = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Erro no servidor", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Completed, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar", http.StatusInternalServerError)
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		http.Error(w, "Tarefa nao encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteTasksHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		http.Error(w, "Erro no servidor", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, "Erro ao deletar", http.StatusInternalServerError)
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		http.Error(w, "Tarefa nao encontrada", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
