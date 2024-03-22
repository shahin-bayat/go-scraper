package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	_ "image/png"

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

func HasImage(base64Img string) (bool, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Img))

	config, _, err := image.DecodeConfig(reader)
	if err != nil {
		return false, err
	}
	return config.Height > 100, nil
}
