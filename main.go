package main

import (
	"github.com/lpernett/godotenv"
	"log"
	"net/http"
	"teachble/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	http.HandleFunc("/", handlers.ServeHome)
	http.HandleFunc("/search", handlers.HandleSearch)
	http.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))

	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
