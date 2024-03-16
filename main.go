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
	var categories []model.Category
	var questions []model.Question
	var _ map[string]string
	var payload request.PreflightPayload
	var err error

	db, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// TODO: DEBUG PURPOSES ONLY - DELETE LATER
	db.Migrator().DropTable(&model.Category{}, &model.Question{}, &model.Answer{})

	db.AutoMigrate(&model.Category{}, &model.Question{}, &model.Answer{})

	initialUrl := util.GetEnvVariable("INITIAL_URL")
	mainUrl := util.GetEnvVariable("MAIN_URL")

	// initial request to fetch cookie and categories
	payload, cookie, err = request.ScrapeInitialPage(initialUrl, db)
	if err != nil {
		panic(err)
	}

	result := db.Model(&model.Category{}).Find(&categories)

	if result.Error != nil {
		log.Fatalf(result.Error.Error())
	}

	// initial POST request defines which federal state's questions to fetch
	// saves question(questionKey and questionNumber) in db
	// from this request on, categoryKey should be questionKey
	// TODO: "4" is the categoryKey, it should be dynamic categoryKey
	payload, err = request.Scrape(mainUrl, cookie, "4", "SUBMIT", payload, db)
	if err != nil {
		panic(err)
	}

	result = db.Model(&model.Question{}).Where("category_id = ? AND is_fetched = ?", 3, false).Limit(3).Find(&questions)
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
