package main

import (
	"fmt"
	"sync"
)

var access = Access{
	htmls: make(map[string]string),
}

type Access struct {
	mutex sync.Mutex
	htmls map[string]string
}

func (a *Access) exists(url string) bool {
	_, ok := a.htmls[url]
	return ok
}

func (a *Access) add(fetcher Fetcher, url string) ([]string, error) {
	if a.exists(url) {
		return nil, nil
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}
	fmt.Printf("found: %s: %s\n", url, body)

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.htmls[url] = body
	return urls, nil
}

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}
	urls, err := access.add(fetcher, url)
	if err != nil {
		fmt.Println(err)
		return
	}

	done := make(chan bool)
	for _, u := range urls {
		go func(url string) {
			Crawl(url, depth-1, fetcher)
			done <- true
		}(u)
	}
	for range urls {
		<-done
	}
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
