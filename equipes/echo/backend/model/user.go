package model

type User struct {
	ID       int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}
