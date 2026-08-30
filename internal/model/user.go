package models

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
