package main

import (
	"fmt"
	"time"
	"tog/fetch"
)

var access = Access{
	htmls: make(map[string]string),
}

type Access struct {
	htmls map[string]string
}

func (a *Access) exists(url string) bool {
	_, ok := a.htmls[url]
	return ok
}

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}
	// Add start
	if access.exists(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s\n", url)
	access.htmls[url] = body
	// Add finish

	for _, u := range urls {
		Crawl(u, depth-1, fetcher)
	}
}

func main() {
	fetcher := fetch.LinkTag{}
	start := time.Now()
	Crawl("https://golang.org/", 2, &fetcher)
	end := time.Now()
	fmt.Printf("%f seconds\n", end.Sub(start).Seconds())
	fmt.Printf("accessed %d htmls\n", len(access.htmls))
}
