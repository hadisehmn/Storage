package auth

import (
	"errors"
	"strings"

	"go-practice/STORAGE/internal/apperror"
	"go-practice/STORAGE/internal/auth/user"
	models "go-practice/STORAGE/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repository *user.UserRepository
}

func NewAuthService(repository *user.UserRepository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (s *AuthService) SignUp(req models.SignUpRequest) (string, error) {

	if strings.TrimSpace(req.Name) == "" {
		return "", errors.New("name is required")
	}

	if strings.TrimSpace(req.Email) == "" {
		return "", errors.New("email is required")
	}

	if req.Password == "" {
		return "", errors.New("password is required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	user := models.User{
		ID:           uuid.New().String(),
		Name:         strings.TrimSpace(req.Name),
		Email:        strings.TrimSpace(req.Email),
		PasswordHash: string(passwordHash),
	}

	if err := s.repository.Create(user); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *AuthService) SignIn(req models.SignInRequest) (string, error) {

	user, err := s.repository.FindByEmail(req.Email)
	if err != nil {
		return "", apperror.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return "", apperror.ErrWrongPassword
	}

	return user.ID, nil
}
