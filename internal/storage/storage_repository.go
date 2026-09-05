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

func (r *StorageRepository) FindByUserID(userID string) ([]models.File, error) {

	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, user_id, file_name, file_path, created_at
		 FROM files
		 WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File

		err := rows.Scan(
			&file.ID,
			&file.UserID,
			&file.FileName,
			&file.FilePath,
			&file.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
