package main

import "time"

type Task struct {
	TaskID    int       `json:"task_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// payload para criação de tarefa
type CreateTaskRequest struct {
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

// payload para atualização de tarefa
type UpdateTaskRequest struct {
	Status string `json:"status"` // "pending" ou "completed"
}