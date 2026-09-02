package main

import (
	"context"
	"log"
	"net/http"

	"go-practice/STORAGE/internal/auth"
	"go-practice/STORAGE/internal/auth/user"
	"go-practice/STORAGE/internal/database"

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

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", authController.SignUp)
	mux.HandleFunc("POST /signin", authController.SignIn)

	log.Println("Server is running on :8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
