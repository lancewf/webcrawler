package webcrawler

import (
	"fmt"
	"github.com/lancewf/concurrent"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var agentStringCollection = concurrent.NewAgentStringCollection([]string{})

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func crawl(url string, depth int, fetcher Fetcher) {
	//println("working", url, depth)
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}

	if containsString(agentStringCollection.GetStringCollection(), url) {
		fmt.Printf("dup1: %s \n", url)
		return
	}

	agentStringCollection.SendStringCollection(func(currentUrls []string) []string {
		//fmt.Printf("in: %s \n", url)
		if containsString(currentUrls, url) {
			return currentUrls
		}

		go fetchAndCrawlLinks(url, depth, fetcher)

		return append(currentUrls, url)
	})

	return
}

func containsString(stringCollection []string, testString string) bool {
	for _, stringInCollection := range stringCollection {
		if stringInCollection == testString {
			return true
		}
	}

	return false
}


func fetchAndCrawlLinks(url string, depth int, fetcher Fetcher) {
	body, foundUrls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range foundUrls {
		//println("starting", u)
		go crawl(u, depth-1, fetcher)
	}
}

func Start() {
	crawl("http://golang.org/", 4, fetcher)

	var input string
	fmt.Scanln(&input)
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
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
