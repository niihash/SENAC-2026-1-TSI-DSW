package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
	"strings"
)

// ================= ROTAS =================

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Método inválido", 405)
		return
	}

	getTasks(w)
}

// ================= FUNCOES PRINCIPAIS =================

// ================= READ =================

func getTasks(w http.ResponseWriter) {
	rows, err := db.Query(`
		SELECT id, user_id, title, done, created_at
		FROM tasks
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Done,
			&t.CreatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		tasks = append(tasks, t)
	}

	json.NewEncoder(w).Encode(tasks)
}