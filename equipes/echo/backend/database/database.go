package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Função auxiliar para ler variáveis do Docker
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Connect() {
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "root")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "todo_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("ERRO: Configuração do banco inválida! ", err)
	}

	// Loop de resiliência. Tenta 5 vezes com pausa de 3 segundos.
	for i := 0; i < 5; i++ {
		err = DB.Ping()
		if err == nil {
			fmt.Println("Conectado ao MySQL com sucesso!")
			return // Sai da função com sucesso
		}
		fmt.Println("Banco ainda não está pronto, aguardando 3 segundos...")
		time.Sleep(3 * time.Second)
	}

	log.Fatal("ERRO FINAL: Não consegui falar com o MySQL! ", err)
}
