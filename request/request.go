package request

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/shahin-bayat/go-scraper/model"
	"gorm.io/gorm"
)

type PreflightPayload struct {
	p_flow_id            string
	p_flow_step_id       string
	p_instance           string
	p_page_submission_id string
	p_request            string
	p_arg_names          string
	p_md5_checksum       string
	p_page_checksum      string
}

func Scrape(url, cookie, questionKey, requestType string, payload PreflightPayload, db *gorm.DB) (PreflightPayload, error) {

	var responsePayload = &PreflightPayload{}

	c := colly.NewCollector()

	var headers map[string]string = map[string]string{
		"User-Agent": UserAgent,
		"Cookie":     cookie,
	}

	c.OnRequest(func(req *colly.Request) {
		SetHeaders(req, &headers)
	})

	// TODO: DEBUG PURPOSES ONLY - DELETE LATER
	c.OnResponse(func(res *colly.Response) {
		fmt.Println("Response received")
		if requestType != "SUBMIT" {
			os.WriteFile(questionKey+".html", res.Body, 0644)
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Fatalf("Ooops, an error occurred!:%s", e)
	})

	// c.OnHTML("table.t3ReportsRegion td.RegionHeader", func(e *colly.HTMLElement) {
	// 	if requestType == "SUBMIT" && strings.Contains(e.Text, "von") {
	// 		questionsCount := strings.Split(e.Text, " ")[3]
	// 		fmt.Println(questionsCount)
	// 	}
	// })

	c.OnHTML("input[type=hidden]", func(e *colly.HTMLElement) {
		name := e.Attr("name")
		value := e.Attr("value")
		switch name {
		case "p_flow_id":
			responsePayload.p_flow_id = value
		case "p_flow_step_id":
			responsePayload.p_flow_step_id = value
		case "p_instance":
			responsePayload.p_instance = value
		case "p_page_submission_id":
			responsePayload.p_page_submission_id = value
		case "p_request":
			responsePayload.p_request = requestType
		case "p_arg_names":
			responsePayload.p_arg_names = value
		case "p_md5_checksum":
			responsePayload.p_md5_checksum = value
		case "p_page_checksum":
			responsePayload.p_page_checksum = value
		}
	})

	c.OnHTML("select[name=p_t01] option", func(e *colly.HTMLElement) {
		if requestType != "SUBMIT" {
			return
		}
		if e.Attr("value") != "%null%" {
			model.SaveQuestion(e.Text, e.Attr("value"), questionKey, db)
		}
	})

	payloadMap := map[string]string{
		"p_flow_id":            payload.p_flow_id,
		"p_flow_step_id":       payload.p_flow_step_id,
		"p_instance":           payload.p_instance,
		"p_page_submission_id": payload.p_page_submission_id,
		"p_request":            payload.p_request,
		"p_arg_names":          payload.p_arg_names,
		"p_t01":                questionKey,
		"p_md5_checksum":       payload.p_md5_checksum,
		"p_page_checksum":      payload.p_page_checksum,
	}

	c.Post(url, payloadMap)

	return *responsePayload, nil

}

func ScrapeInitialPage(url string, db *gorm.DB) (PreflightPayload, string, error) {
	var cookieStr string
	payload := &PreflightPayload{}

	c := colly.NewCollector()

	var headers map[string]string = map[string]string{
		"User-Agent": UserAgent,
	}

	c.OnRequest(func(req *colly.Request) {
		SetHeaders(req, &headers)

	})

	c.OnHTML("input[type=hidden]", func(e *colly.HTMLElement) {
		name := e.Attr("name")
		value := e.Attr("value")
		switch name {
		case "p_flow_id":
			payload.p_flow_id = value
		case "p_flow_step_id":
			payload.p_flow_step_id = value
		case "p_instance":
			payload.p_instance = value
		case "p_page_submission_id":
			payload.p_page_submission_id = value
		case "p_request":
			payload.p_request = "SUBMIT"
		case "p_arg_names":
			payload.p_arg_names = value
		case "p_md5_checksum":
			payload.p_md5_checksum = value
		case "p_page_checksum":
			payload.p_page_checksum = value
		}
	})

	c.OnHTML("select[name=p_t01] option", func(e *colly.HTMLElement) {
		modelCategory := model.Category{
			Text:        e.Text,
			CategoryKey: e.Attr("value"),
		}
		db.Create(&modelCategory)
	})

	c.OnResponse(func(res *colly.Response) {
		cookieStr = strings.Split(res.Headers.Get("Set-Cookie"), ";")[0]
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Fatalf("Blimey, an error occurred!:%s", e)
	})

	c.Visit(url)
	return *payload, cookieStr, nil
}
