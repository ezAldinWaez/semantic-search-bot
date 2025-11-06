package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken   string
	DatabasePath    string
	EmbeddingAPIURL string
	EmbeddingModel  string
	MaxResults      int
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
		log.Println("Using environment variables or defaults")
	}

	return &Config{
		TelegramToken:   getEnv("TELEGRAM_TOKEN", ""),
		DatabasePath:    getEnv("DATABASE_PATH", "./messages.db"),
		EmbeddingAPIURL: getEnv("EMBEDDING_API_URL", "http://localhost:11434"),
		EmbeddingModel:  getEnv("EMBEDDING_MODEL", "all-minilm:latest"),
		MaxResults:      3,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
