package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"go-practice/STORAGE/internal/apperror"
	models "go-practice/STORAGE/internal/model"
)

type AuthController struct {
	service *AuthService
}

func NewAuthController(service *AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (c *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {

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
	w.WriteHeader(http.StatusOK)

	response := models.SignInResponse{
		Message: "login successful",
		Token:   token,
	}

	json.NewEncoder(w).Encode(response)
}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {

	var req models.SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := c.service.SignUp(req)
	if err != nil {
		log.Printf("SignUp failed: %v", err)

		http.Error(
			w,
			"Signup failed",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "signup successful",
		"user_id": userID,
	})
}
