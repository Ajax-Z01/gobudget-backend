package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file (if available)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
	}

	// Initialize database connection and run migrations
	InitDatabase()

	// Seed the database with initial data (only in development mode)
	SeedDatabase()

	// Set up the HTTP router
	router := SetupRouter()

	// Get the server port from environment variables (default to 8080 if not set)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server and listen on the specified port
	log.Println("Server running on port", port)
	router.Run(":" + port)
}
