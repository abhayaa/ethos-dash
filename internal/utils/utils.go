package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvKey(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	return os.Getenv(key)
}
