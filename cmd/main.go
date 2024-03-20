package main

import (
	"errors"
	"log"
	"time"

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

	// STEP 1. Fetch cookie and set categories
	payload, cookie, err = request.ScrapeInitialPage(initialUrl, store)
	if err != nil {
		panic(err)
	}

	// STEP 2: Get category
	category, err := store.GetCategoryByText("Hamburg")
	if err != nil && !errors.Is(err, store.ErrNoRows) {
		log.Fatalf("Error getting category: %s", err.Error())
	}

	// STEP 3: Set questions of the category (key and number)
	payload, err = request.Scrape(mainUrl, cookie, category.CategoryKey, "SUBMIT", payload, store)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// STEP 4: Fix the issue that always the first response is the first question of category, i.e. 17412
	// INFO: This  request won't update the db - look at the logic inside the function
	payload, err = request.Scrape(mainUrl, cookie, "17412", "P30_ROWNUM", payload, store)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// STEP 5: Get questions
	questions, err := store.GetQuestionsByCategoryId(category.ID)
	if err != nil && !errors.Is(err, store.ErrNoRows) {
		log.Fatalf(err.Error())
	}

	// STEP 6: loop through questions and set answers
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
	// delay := time.Duration(util.GenerateRandomDelay(1500, 3000)) * time.Millisecond
	// time.Sleep(delay)
	// 	question := questions[i]
	// 	payload, err = request.Scrape(mainUrl, cookie, question.QuestionKey, "P30_ROWNUM", payload, store)
	// 	if err != nil {
	// 		log.Fatalf(err.Error())
	// 	}
	// }

}
