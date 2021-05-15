package utils

import (
	"os"

	"github.com/joho/godotenv"
)

// GetEnvVar : Load a environment variable from .env file, or the OS
func GetEnvVar(key string) string {
	godotenv.Load("./.env")
	return os.Getenv(key)
}
