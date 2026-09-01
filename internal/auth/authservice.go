package auth

import (
	"go-practice/STORAGE/internal/apperror"
	models "go-practice/STORAGE/internal/model"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) SignIn(req models.SignInRequest) (string, error) {

	if req.Email != "test@gmail.com" {
		return "", apperror.ErrUserNotFound
	}

	if req.Password != "123456" {
		return "", apperror.ErrWrongPassword
	}

	return "user-123", nil
}

func (s *AuthService) SignUp(req models.SignUpRequest) (string, error) {

	return "user-123", nil
}
