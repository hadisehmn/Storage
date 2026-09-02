package storage

type StorageService struct {
	repository *StorageRepository
}

func NewStorageService(repository *StorageRepository) *StorageService {
	return &StorageService{
		repository: repository,
	}
}
