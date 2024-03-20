package main

import (
	"errors"
	"log"
	"time"

	// "time"

	"github.com/shahin-bayat/go-scraper/internal/request"
	"github.com/shahin-bayat/go-scraper/internal/store"
	"github.com/shahin-bayat/go-scraper/internal/util"
)

func main() {
	var cookie string
	var payload request.PreflightPayload
	var err error

	var initialUrl = util.GetEnvVariable("INITIAL_URL")
	var mainUrl = util.GetEnvVariable("MAIN_URL")

	// Init DB
	store, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err = store.Init(); err != nil {
		log.Fatalf(err.Error())
	}

	// STEP 1. Fetch cookie and Categories
	payload, cookie, err = request.ScrapeInitialPage(initialUrl, store)
	if err != nil {
		panic(err)
	}

	// TODO: for debugging purposes only, delete later
	category, err := store.GetCategoryByText("Bayern")
	if err != nil && !errors.Is(err, store.ErrNoRows) {
		log.Fatalf(err.Error())
	}

	// STEP 2: Fetch questions (key and number)
	payload, err = request.Scrape(mainUrl, cookie, category.CategoryKey, "SUBMIT", payload, store)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// TODO: loop through categories and fetch questions
	// categories, err := store.GetCategories()
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// fmt.Println(categories)

	// STEP 3: Loop through questions and fetch answers
	questions, err := store.GetQuestionsByCategoryId(category.ID)
	if err != nil && !errors.Is(err, store.ErrNoRows) {
		log.Fatalf(err.Error())
	}
	for i := 0; i <= 1; i++ {
		delay := time.Duration(util.GenerateRandomDelay(1500, 3000)) * time.Millisecond
		time.Sleep(delay)
		question := questions[i]
		payload, err = request.Scrape(mainUrl, cookie, question.QuestionKey, "P30_ROWNUM", payload, store)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// for i := range questions {
	// 	time.Sleep(3 * time.Second)

	// 	question := questions[i]
	// 	payload, err = request.Scrape(mainUrl, cookie, question.QuestionKey, "P30_ROWNUM", payload, store)
	// 	if err != nil {
	// 		log.Fatalf(err.Error())
	// 	}
	// }

}
