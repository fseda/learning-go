package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type UrlCache struct {
	v map[string]string
	mu sync.Mutex
}

func (c *UrlCache) Set(key string, body string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v[key] = body
} 

func (c *UrlCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.v[key]
	return val, ok
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, urlCache *UrlCache) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	defer wg.Done()

	// Url already fetched
	if _, ok := urlCache.Get(url); ok {
		return
	}

	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	urlCache.Set(url, body)

	// fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, urlCache)
	}
	return
}

var wg sync.WaitGroup

func main() {
	urlCache := &UrlCache{v: make(map[string]string)}

	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher, urlCache)
	wg.Wait()

	for url := range urlCache.v {
		fmt.Printf("%s: %s\n", url, urlCache.v[url])
	}
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
