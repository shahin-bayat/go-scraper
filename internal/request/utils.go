package request

import "github.com/gocolly/colly"

func SetHeaders(req *colly.Request, headers *map[string]string) {
	for k, v := range *headers {
		req.Headers.Set(k, v)
	}
}
