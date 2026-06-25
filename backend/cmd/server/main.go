package main

import (
	"log"

	"emergency-rescue-locator/internal/config"
	"emergency-rescue-locator/internal/database"
	"emergency-rescue-locator/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db := database.Connect(cfg.DatabaseURL)
	router := routes.Setup(db, cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
