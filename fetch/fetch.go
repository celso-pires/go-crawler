package fetch

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LinkTag struct {
}

func (r *LinkTag) Fetch(url string) (string, []string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", nil, fmt.Errorf("http.Get() error: %w", err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("goquery.Get() error: %w", err)
	}
	var urls []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		v, ok := s.Attr("href")
		if ok {
			if strings.Index(v, "/") == 0 {
				v = url + v[1:]
			}
			urls = append(urls, v)
		}
	})

	h, err := doc.Html()
	if err != nil {
		return "", nil, fmt.Errorf("doc.Html() error: %w", err)
	}
	return h, urls, nil
}
