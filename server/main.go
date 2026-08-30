package main

import (
	"log"
	"net/http"

	"go-practice/STORAGE/internal/auth"
	"go-practice/STORAGE/internal/auth/services"
)

func main() {

	userService := services.NewUserService()
	userController := auth.NewUserController(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signin", userController.SignIn)

	log.Println("Server is running on :8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
