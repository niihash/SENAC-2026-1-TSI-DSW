package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initDB()

	r := mux.NewRouter()

	r.HandleFunc("/tasks", getTasksHandler).Methods("GET")

	http.ListenAndServe(":8080", enableCORS(r))
}