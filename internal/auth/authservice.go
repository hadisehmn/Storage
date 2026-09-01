package auth

import (
	"errors"
	"strings"

	"go-practice/STORAGE/internal/apperror"
	models "go-practice/STORAGE/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users map[string]models.User
}

func NewAuthService() *AuthService {
	return &AuthService{
		users: make(map[string]models.User),
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

	for _, user := range s.users {
		if user.Email == req.Email {
			return "", errors.New("email already exists")
		}
	}

	userID := "user-" + string(rune(len(s.users)+1+'0'))

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	user := models.User{
		ID:           userID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}
	s.users[userID] = user

	return userID, nil
}

func (s *AuthService) SignIn(req models.SignInRequest) (string, error) {

	for _, user := range s.users {

		if user.Email == req.Email {

			err := bcrypt.CompareHashAndPassword(
				[]byte(user.PasswordHash),
				[]byte(req.Password),
			)

			if err != nil {
				return "", apperror.ErrWrongPassword
			}

			return user.ID, nil
		}
	}

	return "", apperror.ErrUserNotFound
}
