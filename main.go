package main

import (
	"log"

	"github.com/shahin-bayat/go-scraper/request"
)

const (
	initialUrl string = "https://oet.bamf.de/ords/oetut/f?p=514:1"
	mainUrl    string = "https://oet.bamf.de/ords/oetut/f?p=514:30:0::NO:::"
)

var isInitial bool = true

func main() {
	data, cookieStr, err := request.Scrape(initialUrl, "", "initial.html", &isInitial)
	if err != nil {
		log.Fatalf("Error performing preflight request: %s", err)
	}

	request.Preflight(cookieStr, &data, &isInitial)
	isInitial = false

	request.Scrape(mainUrl, cookieStr, "index.html", &isInitial)

}
