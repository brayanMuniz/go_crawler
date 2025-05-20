# Go Web Crawler

A simple, concurrent web crawler written in Go.  
It crawls a given website, following only internal links, and reports how many times each page was linked to from within the site.

---

## Features

- Concurrency-limited crawling: Control how many pages are fetched in parallel.
- Internal links only: Only follows links within the same host as the starting URL.
- Page visit limit: Stops crawling after a specified number of unique pages.
- Link count report: Outputs a sorted report of how many times each internal page was linked to.

---

## Usage

To run the crawler, use:

    go run . <baseURL> <maxConcurrency> <maxPages>

Where:

- `<baseURL>` is the starting URL to crawl (for example, https://example.com)
- `<maxConcurrency>` is the maximum number of concurrent HTTP requests (for example, 5)
- `<maxPages>` is the maximum number of unique pages to visit (for example, 20)

Example:

    go run . "https://wagslane.dev" 10 25

---

## Output

After crawling, the program prints a report like:

=============================
REPORT for https://example.com
=============================
Found 3 internal links to https://example.com/page1  
Found 2 internal links to https://example.com/page2  
...

Pages are sorted by the number of internal links (descending), then alphabetically.

---

## How It Works

- The crawler starts at the given baseURL.
- It fetches each page, extracts all internal links, and schedules them for crawling (up to the maxPages limit).
- Concurrency is controlled using a semaphore pattern (a buffered channel).
- Only pages with the same host as the baseURL are crawled.
- Each page is counted only once, but the report shows how many times it was linked to from other internal pages.

---

## Design Notes

- **Concurrency:** Uses goroutines and a buffered channel to limit concurrent HTTP requests.
- **Synchronization:** Uses a mutex to protect shared state (the pages map).
- **URL Normalization:** Ensures URLs are compared in a consistent format.
- **Error Handling:** Skips pages that can’t be fetched or parsed, or that aren’t HTML.
- **Sorting:** The report is sorted by link count (descending), then by URL (ascending).

---

## Limitations

- Only follows links with the same host as the starting URL.
- Does not obey robots.txt.
- Does not handle JavaScript-generated links.
- Only processes pages with Content-Type: text/html.
- No retry logic for failed requests.
