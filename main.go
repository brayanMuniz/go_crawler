package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		println("Not enough arguments provided")
		fmt.Println("usage: crawler <baseURL> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}

	if len(os.Args) > 4 {
		println("too many arguments provided")
		fmt.Println("usage: crawler <baseURL> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}

	baseUrl := os.Args[1]
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		println("max concurrency arg is not a number")
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		println("max pages arg is not a number")
		os.Exit(1)
	}

	cfg, err := configure(baseUrl, maxConcurrency, maxPages)

	println("Max concurrency: ", maxConcurrency)
	println("Max pages: ", cfg.maxPages)

	println("starting crawl of: ", baseUrl)

	cfg.wg.Add(1)
	go cfg.crawlPage(baseUrl)
	cfg.wg.Wait()

	printReport(cfg.pages, baseUrl)
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if cfg.pagesLen() >= cfg.maxPages {
		return
	}

	u, err := url.Parse(rawCurrentURL)
	if err != nil {
		print(err)
		return
	}

	resolved := cfg.baseURL.ResolveReference(u)
	if resolved.Host != cfg.baseURL.Host {
		return
	}

	normalURL, err := normalizeURL(resolved.String())
	if err != nil {
		print(err)
		return
	}

	cfg.mu.Lock()
	if cfg.pages[normalURL] != 0 {
		cfg.mu.Unlock()
		return
	} else {
		cfg.pages[normalURL] = 1
		if len(cfg.pages) >= cfg.maxPages {
			println("Hit max amount of pages", len(cfg.pages))
			cfg.mu.Unlock()
			return
		}
	}
	cfg.mu.Unlock()

	extractedHTML, err := getHTML(resolved.String())
	if err != nil {
		return
	}

	// associated urls
	extractedURLS, err := getURLsFromHTML(extractedHTML, cfg.baseURL.String())
	if err != nil {
		return
	} else {
		for _, nextURL := range extractedURLS {
			cfg.wg.Add(1)
			go cfg.crawlPage(nextURL)
		}
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

type PageCount struct {
	URL   string
	Count int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("=============================")

	var visited []PageCount
	for url, count := range pages {
		visited = append(visited, PageCount{URL: url, Count: count})
	}

	sort.Slice(visited, func(i, j int) bool {
		if visited[i].Count != visited[j].Count {
			return visited[i].Count > visited[j].Count // Descending count
		}
		return visited[i].URL < visited[j].URL // Ascending URL
	})

	for _, page := range visited {
		fmt.Printf("Found %d internal links to %s\n", page.Count, page.URL)
	}

}
