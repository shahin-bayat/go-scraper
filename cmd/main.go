package main

import (
	"log"
	"time"

	// "time"

	"github.com/shahin-bayat/go-scraper/internal/request"
	"github.com/shahin-bayat/go-scraper/internal/store"
	"github.com/shahin-bayat/go-scraper/internal/util"
)

func main() {
	var cookie string
	// var categories []model.Category
	// var questions []model.Question
	var _ map[string]string
	var payload request.PreflightPayload
	var err error

	store, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err = store.Init(); err != nil {
		log.Fatalf(err.Error())
	}

	initialUrl := util.GetEnvVariable("INITIAL_URL")
	mainUrl := util.GetEnvVariable("MAIN_URL")

	// STEP 1. Fetch cookie and Categories
	payload, cookie, err = request.ScrapeInitialPage(initialUrl, store)
	if err != nil {
		panic(err)
	}

	// TODO: for debugging purposes only, delete later
	category, err := store.GetCategoryByText("Bayern")
	if err != nil {
		log.Fatalf(err.Error())
	}
	categoryKey := category.CategoryKey

	// STEP 2: Fetch questions (key and number)
	payload, err = request.Scrape(mainUrl, cookie, categoryKey, "SUBMIT", payload, store)
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
	if err != nil {
		log.Fatalf(err.Error())
	}

	for i := 0; i <= 1; i++ {
		time.Sleep(3 * time.Second)
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
