package controller

import (
	"backend/model"
	"backend/repository"
	"encoding/json"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var dadosDigitados model.User

	// Pegamos o JSON enviado no fetch e transformamos em objeto para manipulacao.
	json.NewDecoder(r.Body).Decode(&dadosDigitados)

	// Consultamos o Repository para verificar se esse e-mail existe no banco.
	usuarioDoBanco, _ := repository.GetUserByEmail(dadosDigitados.Email)

	// Comparamos as senhas; se forem diferentes, encerramos com erro de autorização.
	if usuarioDoBanco.Password != dadosDigitados.Password {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Login bem-sucedido: devolvemos os dados do usuário para serem salvos na Promise/Session.
	json.NewEncoder(w).Encode(usuarioDoBanco)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var novoUsuario model.User

	// Traduzimos o JSON recebido para a struct novoUsuario.
	json.NewDecoder(r.Body).Decode(&novoUsuario)

	// Antes de criar, verificamos se o e-mail já consta no banco para evitar duplicidade.
	usuarioExistente, _ := repository.GetUserByEmail(novoUsuario.Email)

	if usuarioExistente.ID > 0 {
		// Retornamos um erro de conflito (409) caso o e-mail já esteja em uso.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "E-mail já cadastrado no sistema"})
		return
	}

	// Persistimos o novo usuário no banco de dados através do Repository.
	repository.CreateUser(novoUsuario)

	// Cadastro finalizado com sucesso.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuário registrado com sucesso"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// apenas avisamos o front para limpar tudo
	json.NewEncoder(w).Encode(model.Response{
		Status:  "success",
		Message: "Logout realizado com sucesso.",
	})
}
