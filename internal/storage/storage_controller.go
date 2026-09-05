package storage

import (
	"encoding/json"
	"net/http"

	"go-practice/STORAGE/internal/auth"
	models "go-practice/STORAGE/internal/model"
)

type StorageController struct {
	service *StorageService
}

func NewStorageController(service *StorageService) *StorageController {
	return &StorageController{
		service: service,
	}
}

const maxFileSize = 100 << 20

func (c *StorageController) Upload(w http.ResponseWriter, r *http.Request) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(
			w,
			"Unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxFileSize,
	)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(
			w,
			"file missing",
			http.StatusBadRequest,
		)
		return
	}

	defer file.Close()

	_ = userID
	_ = header
	_ = file

	err = c.service.Upload(userID, file, header)

	if err != nil {
		http.Error(
			w,
			"upload failed",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("file received"))
}

func (c *StorageController) List(w http.ResponseWriter, r *http.Request) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	files, err := c.service.GetUserFiles(userID)
	if err != nil {
		http.Error(w, "failed to get files", http.StatusInternalServerError)
		return
	}

	response := make([]models.FileResponse, 0, len(files))

	for _, file := range files {
		response = append(response, models.FileResponse{
			ID:        file.ID,
			FileName:  file.FileName,
			CreatedAt: file.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
