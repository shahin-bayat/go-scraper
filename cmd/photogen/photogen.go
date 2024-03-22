package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shahin-bayat/go-scraper/internal/request"
	"github.com/shahin-bayat/go-scraper/internal/store"
	"github.com/shahin-bayat/go-scraper/internal/util"
)

func main() {
	store, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err = store.Init(); err != nil {
		log.Fatalf(err.Error())
	}

	imageBaseUrl := util.GetEnvVariable("IMAGE_BASE_URL")
	questions, err := store.GetQuestions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, question := range questions {
		delay := time.Duration(util.GenerateRandomDelay(1500, 3000)) * time.Millisecond
		time.Sleep(delay)
		request.SaveImage(imageBaseUrl+question.ImagePath, question.QuestionKey, question.ID, store)
		fmt.Printf("Question key %v fetched\n", question.QuestionKey)
	}

	store.Close()
}
