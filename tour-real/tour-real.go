package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"tog/fetch"
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
	fmt.Printf("found: %s\n", url)
	// fmt.Printf("found: %s: %s\n", url, body)

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
	numCPU := runtime.NumCPU()
	fmt.Println("Number of CPU cores:", numCPU)
	fetcher := fetch.LinkTag{}
	start := time.Now()
	Crawl("https://golang.org/", 2, &fetcher)
	end := time.Now()
	fmt.Printf("%f seconds\n", end.Sub(start).Seconds())
	fmt.Printf("accessed %d htmls\n", len(access.htmls))
}
