package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"portone/Routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	routes := routes.SetupRoutes()

	http.Handle("/", routes)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

