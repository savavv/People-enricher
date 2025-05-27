package main

import (
	"log"
	"os"
	"people-enricher/database"
	"people-enricher/routes"

	_ "people-enricher/docs" // swagger docs

	"github.com/joho/godotenv"
)

// @title People Enricher API
// @version 1.0
// @description This is a REST API for enriching people's data.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("[WARN] No .env file found")
	}

	database.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := routes.SetupRouter()

	log.Printf("[INFO] Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[ERROR] Failed to run server: %v", err)
	}
}
