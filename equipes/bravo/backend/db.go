package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func newDB() *sql.DB {
	// Pega as configurações do ambiente ou usa o padrão
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "todolist")

	// Monta a string de conexão
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)

	// Abre a conexão
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("❌ Erro ao abrir conexão: %v", err)
		return db
	}

	// Tenta dar um ping para testar se o banco responde
	if err := db.Ping(); err != nil {
		// MANTEMOS O SERVIDOR VIVO: Trocamos Fatalf por Printf
		log.Printf("⚠️ AVISO: Banco de dados não encontrado ou desligado (%v).", err)
		log.Println("🚀 O servidor continuará rodando para testes de integração de API.")
	} else {
		log.Println("✅ Conexão com o banco estabelecida com sucesso!")
	}

	return db
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}