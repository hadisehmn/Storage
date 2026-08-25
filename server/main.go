package main

import (
	"go-practice/STORAGE/internal/user"
	"log"
	"net/http"
)

func main() {

	storageController := &user.StorageController{}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /storage/upload", storageController.Upload)

	log.Println("Server is running on :8080")
	http.ListenAndServe(":8080", mux)
}
