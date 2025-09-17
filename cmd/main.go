package main

import (
	"log"
	"os"

	"product-service/internal/httpserver"
)

func main() {
	server, err := httpserver.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Product service starting on port %s", port)
	
	if err := server.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}