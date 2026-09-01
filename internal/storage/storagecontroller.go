package storage

import (
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

	// 1. گرفتن userID از context
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(
			w,
			"Unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	// 2. محدود کردن حجم Request
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxFileSize,
	)

	// 3. گرفتن فایل از HTTP Request
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

	// 4. فعلاً فقط برای اینکه ببینیم اطلاعات درست گرفته شده
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

func (c *StorageController) Get(w http.ResponseWriter, r *http.Request) {
	// HTTP مربوط به Get
}

func (c *StorageController) Delete(w http.ResponseWriter, r *http.Request) {
	// HTTP مربوط به Delete
}
