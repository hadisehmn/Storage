package storage

import models "go-practice/STORAGE/internal/model"

func (s *StorageService) DownloadFile(
	fileID string,
	userID string,
) (*models.File, error) {

	return s.repository.FindByIDAndUserID(fileID, userID)
}
