package user

import (
	"fmt"
	"io"
	"net/http"
)

type StorageController struct{}

const maxFileSize = 100 << 20

func (h *StorageController) Upload(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > maxFileSize {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "cannot read file", http.StatusBadRequest)
		return
	}

	mimeType := http.DetectContentType(buffer)
	fmt.Println("MIME:", mimeType)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, "cannot reset file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("file received"))
}
