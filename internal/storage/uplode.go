package storage

import "mime/multipart"

func (s *StorageService) Upload(
	userID string,
	file multipart.File,
	header *multipart.FileHeader,
) error {

	// منطق واقعی Upload

	return nil
}
