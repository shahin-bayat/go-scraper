package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
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

func ConvertToBase64(data []byte, contentType string) (string, error) {
	data, err := io.ReadAll(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	base64Encoded := base64.StdEncoding.EncodeToString(data)
	formattedBase64 := fmt.Sprintf("data:%s;base64,%s", contentType, base64Encoded)
	return formattedBase64, nil

}
