package storage

import (
	"encoding/json"
	"go-practice/STORAGE/internal/auth"
	"net/http"
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
