package user

import (
	"context"

	"go-practice/STORAGE/internal/model"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user model.User) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO users
			(id, name, email, password_hash)
		 VALUES ($1, $2, $3, $4)`,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, name, email, password_hash
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
