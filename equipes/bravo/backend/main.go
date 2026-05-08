package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	// Import para resolver o problema de conexão entre Front e Back (CORS)
	// IMPORTANTE: Se der erro de "package not found", rode: go get github.com/rs/cors
	"github.com/rs/cors"
)

type App struct {
	db *sql.DB
}

func main() {
	// 1. Inicializa a conexão com o banco de dados
	db := newDB()
	defer db.Close()

	app := &App{db: db}

	// 2. Define o Mux (Roteador)
	mux := http.NewServeMux()

	// 3. Registra as rotas de tarefas
	// /api/v1/tasks -> GET (listar) | POST (criar)
	mux.HandleFunc("/api/v1/tasks", app.tasksCollectionHandler)
	// /api/v1/tasks/{id} -> GET | PUT | DELETE
	mux.HandleFunc("/api/v1/tasks/", app.taskItemHandler)

	// 4. Configuração do CORS (O segredo para o JS funcionar)
	// Isso permite que o seu Frontend acesse a API mesmo estando em portas diferentes
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Permite qualquer origem para teste local
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Envolve o roteador com o middleware de CORS
	handler := c.Handler(mux)

	// 5. Inicia o servidor
	port := getEnv("PORT", "8080")
	log.Printf("🚀 Servidor Bravo rodando na porta %s", port)
	
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}

// Handler para a coleção de tarefas (/api/v1/tasks)
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

// Handler para itens individuais (/api/v1/tasks/{id})
func (app *App) taskItemHandler(w http.ResponseWriter, r *http.Request) {
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