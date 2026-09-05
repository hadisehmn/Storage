package main

import (
	"context"
	"log"
	"net/http"

	"go-practice/STORAGE/internal/auth"
	"go-practice/STORAGE/internal/auth/user"
	"go-practice/STORAGE/internal/database"
	"go-practice/STORAGE/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	userRepository := user.NewUserRepository(db)

	authService := auth.NewAuthService(userRepository)
	authController := auth.NewAuthController(authService)

	storageRepository := storage.NewStorageRepository(db)
	storageService := storage.NewStorageService(storageRepository)
	storageController := storage.NewStorageController(storageService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", authController.SignUp)
	mux.HandleFunc("POST /signin", authController.SignIn)

	mux.Handle(
		"POST /upload",
		auth.Authentication(http.HandlerFunc(storageController.Upload)),
	)

	mux.Handle(
		"GET /files",
		auth.Authentication(http.HandlerFunc(storageController.List)),
	)

	log.Println("Server is running on :8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
