package main

import (
	"strconv"
	"time"

	"github.com/shahin-bayat/go-scraper/request"
	"github.com/shahin-bayat/go-scraper/util"
)

func main() {
	var cookie string
	var categories map[string]string
	var pages map[string]string
	var payload request.PreflightPayload
	var err error

	initialUrl := util.GetEnvVariable("INITIAL_URL")
	mainUrl := util.GetEnvVariable("MAIN_URL")

	payload, cookie, categories, err = request.ScrapeInitialPage(initialUrl)
	if err != nil {
		panic(err)
	}

	payload, pages, err = request.Scrape(mainUrl, cookie, categories["Bayern"], "SUBMIT", payload, "index1.html")
	if err != nil {
		panic(err)
	}

	for i := 1; i <= 3; i++ {
		time.Sleep(3 * time.Second)
		payload, _, err = request.Scrape(mainUrl, cookie, pages[strconv.Itoa(i)], "P30_ROWNUM", payload, "index"+strconv.Itoa(i)+".html")
		if err != nil {
			panic(err)
		}
	}

}
