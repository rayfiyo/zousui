package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var GeminiAPIKEY string

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load the .env: %v", err)
		return err
	}

	GeminiAPIKEY = os.Getenv("GEMINI_API_KEY")

	return nil
}
