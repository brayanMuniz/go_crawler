package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

	println("starting crawl of: ", os.Args[1])

	htmlString, err := getHTML(os.Args[1])
	if err != nil {
		println("Error reading html")
		os.Exit(1)
	}
	print(htmlString)
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Error was of status code 400+")
	}

	if resp.Header.Get("Content-Type") != "text/html" {
		return "", fmt.Errorf("Content-Type is not html")
	}

	htmlString, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(htmlString), nil

}
