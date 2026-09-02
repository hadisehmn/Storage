package storage

import (
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/google/uuid"

	models "go-practice/STORAGE/internal/model"
)

func (s *StorageService) Upload(
	userID string,
	file multipart.File,
	header *multipart.FileHeader,
) error {

	fileID := uuid.New().String()

	filePath := "uploads/" + fileID

	err := os.MkdirAll("uploads", 0755)
	if err != nil {
		return err
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	fileRecord := models.File{
		ID:        fileID,
		UserID:    userID,
		FileName:  header.Filename,
		FilePath:  filePath,
		CreatedAt: time.Now(),
	}

	err = s.repository.Create(fileRecord)
	if err != nil {
		os.Remove(filePath)
		return err
	}

	return nil
}
