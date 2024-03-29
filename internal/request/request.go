package request

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/shahin-bayat/go-scraper/internal/model"
	"github.com/shahin-bayat/go-scraper/internal/store"
	"github.com/shahin-bayat/go-scraper/internal/util"
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

func Scrape(url, cookie, questionKey, requestType string, payload PreflightPayload, store *store.Store) (PreflightPayload, error) {
	var responsePayload = &PreflightPayload{}

	c := colly.NewCollector()

	var headers map[string]string = map[string]string{
		"User-Agent": UserAgent,
		"Cookie":     cookie,
	}

	c.OnRequest(func(req *colly.Request) {
		SetHeaders(req, &headers)
	})

	c.OnResponse(func(res *colly.Response) {
		if requestType != "SUBMIT" {
			fmt.Printf("question key %v: response received\n", questionKey)
			os.WriteFile("web/"+questionKey+".html", res.Body, 0644)
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Fatalf("Ooops, an error occurred!:%s", e)
	})

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
			categoryKey := questionKey
			category, err := store.GetCategoryByCategoryKey(categoryKey)
			if err != nil && !errors.Is(err, store.ErrNoRows) {
				log.Fatalf("Error fetching category:%s", err)
			}

			existingQuestion, err := store.GetQuestionByQuestionKey(e.Attr("value"))
			if err != nil && !errors.Is(err, store.ErrNoRows) {
				log.Fatalf("Error fetching question:%s", err)
			}

			if existingQuestion == nil {
				question := model.CreateQuestion(e.Text, e.Attr("value"))
				questionId, err := store.CreateQuestion(question, category)
				if err != nil {
					log.Fatalf("Error creating question:%s", err)
				}
				if err = store.AssociateQuestionWithCategory(category, questionId); err != nil {
					log.Fatalf("Error associating question with category:%s", err)
				}
			} else {
				isAssociated, err := store.IsQuestionAssociatedWithCategory(category.ID, existingQuestion.ID)
				if err != nil && !errors.Is(err, store.ErrNoRows) {
					log.Fatalf("Error fetching question:%s", err)
				}
				if isAssociated {
					return
				}
				if err := store.AssociateQuestionWithCategory(category, existingQuestion.ID); err != nil {
					log.Fatalf("Error associating question with category:%s", err)
				}
			}
		}
	})

	c.OnHTML("span#P30_AUFGABENSTELLUNG_BILD img", func(e *colly.HTMLElement) {
		if requestType != "P30_ROWNUM" {
			return
		}
		question, err := store.GetQuestionByQuestionKey(questionKey)
		if err != nil && !errors.Is(err, store.ErrNoRows) {
			log.Fatalf("Error fetching question:%s", err)
		}
		updatedQuestion := model.UpdateQuestion(question, &model.UpdateQuestionRequest{ImagePath: e.Attr("src")})
		if err := store.UpdateQuestion(question.ID, updatedQuestion); err != nil {
			log.Fatalf("Error updating question:%s", err)
		}
	})

	c.OnHTML("tr td[headers=RICHTIGE_ANTWORT]", func(e *colly.HTMLElement) {
		if requestType != "P30_ROWNUM" {
			return
		}
		answerText := e.DOM.Parent().Find("td[headers=ANTWORT]").Text()
		isCorrect := strings.Contains(e.Text, "richtige Antwort")
		question, err := store.GetQuestionByQuestionKey(questionKey)
		if err != nil && !errors.Is(err, store.ErrNoRows) {
			log.Fatalf("Error fetching question:%s", err)
		}
		isFetched, err := store.IsQuestionFetched(question.ID)
		if err != nil && !errors.Is(err, store.ErrNoRows) {
			log.Fatalf("Error fetching question:%s", err)
		}
		if isFetched {
			return
		}

		answer := model.CreateAnswer(answerText, isCorrect, question)
		if err := store.CreateAnswer(answer); err != nil {
			log.Fatalf("Error creating answer:%s", err)
		}

		// isFetched = true if there are 4 answers
		answers, err := store.GetAnswersByQuestionId(question.ID)
		if err != nil && !errors.Is(err, store.ErrNoRows) {
			log.Fatalf("Error fetching answers:%s", err)
		}
		if len(answers) == 4 {
			updatedQuestion := model.UpdateQuestion(question, &model.UpdateQuestionRequest{IsFetched: true})
			if err := store.UpdateQuestion(question.ID, updatedQuestion); err != nil {
				log.Fatalf("Error updating question:%s", err)
			}
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

func ScrapeInitialPage(url string, store *store.Store) (PreflightPayload, string, error) {
	var cookieStr string
	responsePayload := &PreflightPayload{}

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
			responsePayload.p_flow_id = value
		case "p_flow_step_id":
			responsePayload.p_flow_step_id = value
		case "p_instance":
			responsePayload.p_instance = value
		case "p_page_submission_id":
			responsePayload.p_page_submission_id = value
		case "p_request":
			responsePayload.p_request = "SUBMIT"
		case "p_arg_names":
			responsePayload.p_arg_names = value
		case "p_md5_checksum":
			responsePayload.p_md5_checksum = value
		case "p_page_checksum":
			responsePayload.p_page_checksum = value
		}
	})

	c.OnHTML("select[name=p_t01] option", func(e *colly.HTMLElement) {
		category := model.CreateCategory(e.Text, e.Attr("value"))
		store.CreateCategory(category)
	})

	c.OnResponse(func(res *colly.Response) {
		cookieStr = strings.Split(res.Headers.Get("Set-Cookie"), ";")[0]
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Fatalf("Blimey, an error occurred!:%s", e)
	})

	c.Visit(url)
	return *responsePayload, cookieStr, nil
}

func SaveImage(imageUrl, questionKey string, questionId uint, store *store.Store) (string, error) {
	c := colly.NewCollector()
	var base64Encoded string

	var headers map[string]string = map[string]string{
		"User-Agent": UserAgent,
	}

	c.OnRequest(func(req *colly.Request) {
		SetHeaders(req, &headers)
	})

	c.OnResponse(func(res *colly.Response) {
		var err error
		filename := questionKey + ".png"
		os.WriteFile("assets/images/"+filename, res.Body, 0644)

		contentType := res.Headers.Get("Content-Type")
		base64Encoded, err = util.ConvertToBase64(res.Body, contentType)
		if err != nil {
			log.Fatalf("Error converting image to base64:%s", err)
		}
		os.WriteFile("assets/base64/"+questionKey+".txt", []byte(base64Encoded), 0644)

		image := model.CreateImage(questionId, filename)
		if err := store.CreateImage(image); err != nil {
			log.Fatalf("Error creating image:%s", err)
		}
	})
	c.OnError(func(r *colly.Response, e error) {
		log.Fatalf("Oops, an error occurred!:%s", e)
	})

	c.Visit(imageUrl)

	return base64Encoded, nil
}
