package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"go-practice/STORAGE/internal/apperror"
	"go-practice/STORAGE/internal/auth/services"
	models "go-practice/STORAGE/internal/model"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (c *UserController) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := c.service.SignIn(req)
	if err != nil {
		log.Printf("SignIn failed: %v", err)

		switch {
		case errors.Is(err, apperror.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound)

		case errors.Is(err, apperror.ErrWrongPassword):
			http.Error(w, "Wrong password", http.StatusUnauthorized)

		default:
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
		}

		return
	}

	token, err := GenerateToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "login successful",
		"token":   token,
	})
}
