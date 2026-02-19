package main

import (
	"log"
	"net/http"
	"os"

	"github.com/martinzha28/bridge-platform/backend/internal/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := routes.SetupRouter()

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
