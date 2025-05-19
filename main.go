package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		println("too many arguments provided")
		os.Exit(1)
	}

	baseUrl := os.Args[1]
	println("starting crawl of: ", baseUrl)

	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		fmt.Println("Not able to parse the base url")
		os.Exit(1)
	}

	visited := make(map[string]int)

	htmlString, err := getHTML(baseUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	visited[baseUrlParsed.String()] += 1

	extractedUrls, err := getURLsFromHTML(htmlString, baseUrl)
	if err != nil {
		println("Could not extract urls from html")
		os.Exit(1)
	}

	for _, extracted_url := range extractedUrls {
		crawlPage(baseUrlParsed, extracted_url, visited)
	}

}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Error was of status code 400+")
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("Content-Type is not html")
	}

	htmlString, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(htmlString), nil
}

func crawlPage(baseUrlParsed *url.URL, rawCurrentURL string, pages map[string]int) {
	u, err := url.Parse(rawCurrentURL)
	if err != nil {
		print(err)
		return
	}

	if pages[u.String()] != 0 {
		println("Already visited", u.String())
		return
	}

	resolved := baseUrlParsed.ResolveReference(u)
	if resolved.Host != baseUrlParsed.Host {
		println("Host is not the same: ", u.Host)
		return
	}

	extractedHTML, err := getHTML(u.String())
	if err != nil {
		fmt.Println("Failed to get HTML for", u.String())
		return
	} else {
		pages[rawCurrentURL] += 1
		// get only associated urls
		extractedURLS, err := getURLsFromHTML(extractedHTML, baseUrlParsed.String())
		if err != nil {
			fmt.Println(err)
			return
		} else {
			for _, u := range extractedURLS {
				crawlPage(baseUrlParsed, u, pages)
			}
		}

	}

}
