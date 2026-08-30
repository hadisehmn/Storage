package services

import (
	"go-practice/STORAGE/internal/apperror"
	models "go-practice/STORAGE/internal/model"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) SignIn(req models.SignInRequest) (string, error) {

	if req.Email != "test@test.com" {
		return "", apperror.ErrUserNotFound
	}

	if req.Password != "123456" {
		return "", apperror.ErrWrongPassword
	}

	return "1", nil
}
