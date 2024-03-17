package main

import (
	"log"
	"time"

	"github.com/shahin-bayat/go-scraper/model"
	"github.com/shahin-bayat/go-scraper/request"
	"github.com/shahin-bayat/go-scraper/store"
	"github.com/shahin-bayat/go-scraper/util"
)

func main() {
	var cookie string
	// var categories []model.Category
	var questions []model.Question
	var _ map[string]string
	var payload request.PreflightPayload
	var err error

	db, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// INFO: DEBUG PURPOSES ONLY - DELETE LATER
	// db.Migrator().DropTable(&model.Category{}, &model.Question{}, &model.Answer{})

	db.AutoMigrate(&model.Category{}, &model.Question{}, &model.Answer{})

	initialUrl := util.GetEnvVariable("INITIAL_URL")
	mainUrl := util.GetEnvVariable("MAIN_URL")

	// STEP 1. Fetch cookie and Categories
	payload, cookie, err = request.ScrapeInitialPage(initialUrl, db)
	if err != nil {
		panic(err)
	}

	// FIXME: "3" is the categoryKey, it should be dynamic
	// STEP 2: Fetch questions (key and number)
	payload, err = request.Scrape(mainUrl, cookie, "3", "SUBMIT", payload, db)
	if err != nil {
		panic(err)
	}

	// FIXME: 2 is the categoryID for categoryKey = 3, it should be dynamic
	// STEP 3: Loop through questions and fetch answers
	result := db.Model(&model.Question{}).Where("category_id = ? AND is_fetched = ?", 2, false).Limit(3).Find(&questions)
	if result.Error != nil {
		log.Fatalf(result.Error.Error())
	}

	for i := range questions {
		time.Sleep(3 * time.Second)

		question := questions[i]
		payload, err = request.Scrape(mainUrl, cookie, question.QuestionKey, "P30_ROWNUM", payload, db)
		if err != nil {
			panic(err)
		}
	}

}
