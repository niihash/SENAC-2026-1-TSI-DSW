package repository

import (
	"backend/database"
	"backend/model"
)

// criar novo usuario
func CreateUser(user model.User) error {
	_, err := database.DB.Exec(
		`INSERT INTO users (email, password_hash) VALUES (?, ?)`,
		user.Email,
		user.Password,
	)
	return err
}

// GetUserByEmail: Vai no banco buscar os dados salvos para conferir no Login.
func GetUserByEmail(email string) (model.User, error) {
	var user model.User

	// Executa a busca e já preenche o objeto 'user' com os dados do banco.
	err := database.DB.QueryRow(
		`SELECT id, email, password_hash FROM users WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		// Se não achar nada, devolve um usuário vazio e o aviso de erro.
		return model.User{}, err
	}

	return user, nil
}
