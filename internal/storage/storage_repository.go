package storage

import (
	"context"

	models "go-practice/STORAGE/internal/model"

	"github.com/jackc/pgx/v5"
)

type StorageRepository struct {
	db *pgx.Conn
}

func NewStorageRepository(db *pgx.Conn) *StorageRepository {
	return &StorageRepository{
		db: db,
	}
}

func (r *StorageRepository) Create(file models.File) error {

	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO files
			(id, user_id, file_name, file_path, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		file.ID,
		file.UserID,
		file.FileName,
		file.FilePath,
		file.CreatedAt,
	)

	return err
}
