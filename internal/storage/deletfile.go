package storage

import "os"

func (s *StorageService) DeleteFile(fileID string, userID string) error {

	file, err := s.repository.FindByIDAndUserID(fileID, userID)
	if err != nil {
		return err
	}

	err = os.Remove(file.FilePath)
	if err != nil {
		return err
	}

	return s.repository.Delete(fileID, userID)
}
