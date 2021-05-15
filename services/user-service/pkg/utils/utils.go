package utils

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// GetEnvVar : Load a environment variable from .env file, or the OS
func GetEnvVar(key string) string {
	godotenv.Load("./.env")
	return os.Getenv(key)
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}
