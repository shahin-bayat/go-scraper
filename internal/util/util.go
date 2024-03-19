package util

import (
	"log"
	"math/rand"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func GenerateRandomDelay(minMs uint, maxMs uint) uint {
	randomDelayMs := minMs + uint(rand.Intn(int(maxMs-minMs+1)))
	return randomDelayMs
}
