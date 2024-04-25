package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type Fetcher interface {
	Fetch(url string) ([]string, error)
}

type fetch struct {
	fetched map[string]bool
}

func (f fetch) Fetch(url string, useCache bool) ([]string, error) {
	urls := []string{}
	if useCache {
		url = GOOGLE_CACHE + url
	}
	ctx := context.Background()

	c := client()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "http://www.google.com/")

	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var fn func(*html.Node)
	fn = func(h *html.Node) {
		if h.Type == html.ElementNode && h.Data == "a" {
			for _, a := range h.Attr {
				if a.Key == "href" {
					urls = append(urls, a.Val)
				}
			}
		}

		for c := h.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}

	fn(doc)
	if len(urls) <= 0 {
		log.Printf("found zero urls\n")
	}

	return urls, nil
}
