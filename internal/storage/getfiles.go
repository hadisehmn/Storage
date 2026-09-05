package storage

import models "go-practice/STORAGE/internal/model"

func (s *StorageService) GetUserFiles(userID string) ([]models.File, error) {
	return s.repository.FindByUserID(userID)
}
