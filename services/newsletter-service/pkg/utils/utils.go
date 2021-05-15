package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVar(key string) string {
	godotenv.Load("./.env")
	return os.Getenv(key)
}
