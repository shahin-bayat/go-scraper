package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func BasicScrapper() {
	url := "https://oet.bamf.de/ords/oetut/f?p=514:1"

	// Create a custom HTTP client
	client := &http.Client{}

	// Set headers in the custom request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "oet.bamf.de")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"122\", \"Not(A:Brand\";v=\"24\", \"Google Chrome\";v=\"122\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Cookie", "AODSESSION=5f9d58b86e1138a226e3f299299261944a9ca7af099c6bc3151c7277b97538fd")

	// Perform the request using the custom client and request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("getting %s: %s", url, resp.Status)
		os.Exit(1)
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("parsing %s as HTML: %v", url, err)
		os.Exit(1)
	}

	// Get HTML from parsed node
	var buf bytes.Buffer
	if err := html.Render(&buf, node); err != nil {
		log.Fatalf("rendering HTML: %v", err)
		os.Exit(1)
	}

	htmlContent := buf.String()

	fmt.Println(htmlContent)

}

// func findStringSubmatch(regex string, content string) []string {
// 	re := regexp.MustCompile(regex)
// 	return re.FindStringSubmatch(content)
// }
