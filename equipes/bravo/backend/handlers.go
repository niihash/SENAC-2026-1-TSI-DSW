package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// POST /api/v1/tasks
func (app *App) createTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "o campo 'name' é obrigatório")
		return
	}
	if req.UserID == 0 {
		writeError(w, http.StatusBadRequest, "o campo 'user_id' é obrigatório")
		return
	}

	result, err := app.db.Exec(
		"INSERT INTO tasks (name, user_id) VALUES (?, ?)",
		req.Name, req.UserID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao criar tarefa")
		return
	}

	id, _ := result.LastInsertId()

	var task Task
	err = app.db.QueryRow(
		"SELECT task_id, name, status, user_id, created_at, updated_at FROM tasks WHERE task_id = ?", id,
	).Scan(&task.TaskID, &task.Name, &task.Status, &task.UserID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar tarefa criada")
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

// GET /api/v1/tasks
func (app *App) listTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := app.db.Query(
		"SELECT task_id, name, status, user_id, created_at, updated_at FROM tasks ORDER BY created_at DESC",
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar tarefas")
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.TaskID, &t.Name, &t.Status, &t.UserID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "erro ao ler tarefa")
			return
		}
		tasks = append(tasks, t)
	}

	writeJSON(w, http.StatusOK, tasks)
}

// PUT /api/v1/tasks/{id}
func (app *App) updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := extractIDFromPath(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if req.Status != "pending" && req.Status != "completed" {
		writeError(w, http.StatusBadRequest, "status deve ser 'pending' ou 'completed'")
		return
	}

	result, err := app.db.Exec(
		"UPDATE tasks SET status = ? WHERE task_id = ?",
		req.Status, id,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao atualizar tarefa")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "tarefa não encontrada")
		return
	}

	var task Task
	err = app.db.QueryRow(
		"SELECT task_id, name, status, user_id, created_at, updated_at FROM tasks WHERE task_id = ?", id,
	).Scan(&task.TaskID, &task.Name, &task.Status, &task.UserID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar tarefa atualizada")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// DELETE /api/v1/tasks/{id}
func (app *App) deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := extractIDFromPath(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	result, err := app.db.Exec("DELETE FROM tasks WHERE task_id = ?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao deletar tarefa")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "tarefa não encontrada")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/tasks/{id}
func (app *App) getTask(w http.ResponseWriter, r *http.Request) {
	idStr := extractIDFromPath(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	var task Task
	err = app.db.QueryRow(
		"SELECT task_id, name, status, user_id, created_at, updated_at FROM tasks WHERE task_id = ?", id,
	).Scan(&task.TaskID, &task.Name, &task.Status, &task.UserID, &task.CreatedAt, &task.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "tarefa não encontrada")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar tarefa")
		return
	}

	writeJSON(w, http.StatusOK, task)
}