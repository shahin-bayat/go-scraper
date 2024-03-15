package request

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type PreflightRequest struct {
	p_flow_id            string
	p_flow_step_id       string
	p_instance           string
	p_page_submission_id string
	p_request            string
	p_arg_names          string
	p_t01                string
	p_md5_checksum       string
	p_page_checksum      string
}

func Scrape(url string, cookie string, filename string, isInitial *bool) (PreflightRequest, string, error) {
	data := &PreflightRequest{}
	var cookieStr string

	c := colly.NewCollector()

	// set http request headers
	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		if cookie != "" {
			req.Headers.Set("Cookie", cookie)
		}
	})

	c.OnResponse(func(res *colly.Response) {
		fmt.Println("Response received")
		if *isInitial {
			cookieStr = strings.Split(res.Headers.Get("Set-Cookie"), ";")[0]
		}
		os.WriteFile(filename, res.Body, 0644)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Fatalf("Blimey, an error occurred!:%s", e)
	})

	c.OnHTML("input[type=hidden]", func(e *colly.HTMLElement) {
		name := e.Attr("name")
		value := e.Attr("value")
		switch name {
		case "p_flow_id":
			data.p_flow_id = value
		case "p_flow_step_id":
			data.p_flow_step_id = value
		case "p_instance":
			data.p_instance = value
		case "p_page_submission_id":
			data.p_page_submission_id = value
		case "p_request":
			if *isInitial {
				data.p_request = "SUBMIT"
			} else {
				data.p_request = "P30_ROWNUM"
			}
		case "p_arg_names":
			data.p_arg_names = value
		case "p_md5_checksum":
			if value != "" {
				data.p_md5_checksum = value
			} else {
				data.p_md5_checksum = ""
			}
		case "p_page_checksum":
			data.p_page_checksum = value
		}
		data.p_t01 = "4" // TODO: first request defines the state, the next ones are question numbers from select
	})

	error := c.Visit(url)
	if error != nil {
		log.Fatal(error)
	}
	return *data, cookieStr, nil

}

func Preflight(cookie string, preflightData *PreflightRequest, initialRequest *bool) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://oet.bamf.de/ords/oetut/wwv_flow.accept", nil)
	if err != nil {
		log.Fatalf("Error creating preflight request: %s", err)
	}

	var referer string
	if *initialRequest {
		referer = FirstReferer
	} else {
		referer = Referer
	}

	var headers map[string]string = map[string]string{
		"Accept":          Accept,
		"Accept-Encoding": AcceptEncoding,
		"Accept-Language": AcceptLanguage,
		"Cache-Control":   CacheControl,
		"Connection":      Connection,
		"Content-Length":  "211",
		"Content-Type":    ContentType,
		"Cookie":          cookie,
		"Host":            Host,
		"Origin":          Origin,
		"Referer":         referer,
		"User-Agent":      UserAgent,
	}

	setHeaders(req, &headers)

	// set body form data
	form := url.Values{}
	form.Add("p_flow_id", preflightData.p_flow_id)
	form.Add("p_flow_step_id", preflightData.p_flow_step_id)
	form.Add("p_instance", preflightData.p_instance)
	form.Add("p_page_submission_id", preflightData.p_page_submission_id)
	form.Add("p_request", preflightData.p_request)
	form.Add("p_arg_names", preflightData.p_arg_names)
	form.Add("p_t01", preflightData.p_t01) // question number in options
	form.Add("p_md5_checksum", preflightData.p_md5_checksum)
	form.Add("p_page_checksum", preflightData.p_page_checksum)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing preflight request: %s", err)
	}

	fmt.Println(resp.Status)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading preflight response: %s", err)
	}

	os.WriteFile("preflight.html", body, 0644)
}

func setHeaders(req *http.Request, headers *map[string]string) {
	for k, v := range *headers {
		req.Header.Set(k, v)
	}
}
