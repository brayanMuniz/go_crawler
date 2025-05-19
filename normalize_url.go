package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL: %w", err)
	}

	fullPath := parsedURL.Host + parsedURL.Path

	fullPath = strings.ToLower(fullPath)

	fullPath = strings.TrimSuffix(fullPath, "/")

	return fullPath, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseUrl, err := url.Parse(rawBaseURL)
	if err != nil {
		return make([]string, 0), nil
	}

	allUrls := make([]string, 0)
	htmlReader := strings.NewReader(htmlBody)

	node, err := html.Parse(htmlReader)
	if err != nil {
		return nil, err
	}

	var visit func(*html.Node)
	visit = func(n *html.Node) {
		if n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := url.Parse(attr.Val)
					if err != nil {
						return
					}

					if u.IsAbs() {
						allUrls = append(allUrls, u.String())
					} else {
						resolved := baseUrl.ResolveReference(u)
						allUrls = append(allUrls, resolved.String())
					}

				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(node)

	return allUrls, nil
}
