package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
)

// WebCrawler contains all the data needed to crawl the web URL
type WebCrawler struct {
	client              *http.Client
	initialURL          string
	nWorkers            int
	outputLinks         []string
	crawledLinksChannel chan string
	pendingURLChannel   chan string
	pendingCountChannel chan int
	mx                  sync.Mutex
}

// NewWebCrawler return a new WebCrawler with a default/custom config
func NewWebCrawler() *WebCrawler {
	var url string
	var nWorkers int
	flag.StringVar(&url, "url", "https://parserdigital.com/", "First URL to Crawl")
	flag.IntVar(&nWorkers, "n", 10, "Number of max workers")

	flag.Parse()
	return &WebCrawler{
		client:              initClient(),
		initialURL:          url,
		nWorkers:            nWorkers,
		outputLinks:         []string{},
		crawledLinksChannel: make(chan string),
		pendingURLChannel:   make(chan string),
		pendingCountChannel: make(chan int),
	}
}

// Run start crawling
func (cwl *WebCrawler) Run() {

	go func() {
		cwl.crawledLinksChannel <- crawler.initialURL
	}()

	var wg sync.WaitGroup
	go cwl.LinkHandler()
	go func() {
		if cwl.CountDown() {
			close(cwl.pendingURLChannel)
			close(cwl.crawledLinksChannel)
			close(cwl.pendingCountChannel)
		}
	}()

	for w := 1; w <= crawler.nWorkers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			defer crawler.mx.Unlock()

			links := cwl.CrawlLinks()

			crawler.mx.Lock()
			crawler.outputLinks = append(crawler.outputLinks, links...)
		}(w)

	}

	wg.Wait()
	fmt.Printf("%dx Workers found a total of %d valid Links:\n%+v\n\n", crawler.nWorkers, len(crawler.outputLinks), crawler.outputLinks)
}

// CrawlLinks will inspect the URL from the pending URL channel
func (cwl *WebCrawler) CrawlLinks() []string {
	links := []string{}
	for link := range cwl.pendingURLChannel {
		inspectURLContent(cwl, link)
		cwl.pendingCountChannel <- -1
		links = append(links, link)
	}
	return links
}

// LinkHandler filter the crawled link and marked them as Read
func (cwl *WebCrawler) LinkHandler() {
	alreadyCrawled := make(map[string]bool)

	for link := range cwl.crawledLinksChannel {
		if !alreadyCrawled[link] {
			// mark as crawled
			alreadyCrawled[link] = true
			// send it to start crawling it
			cwl.pendingCountChannel <- 1
			cwl.pendingURLChannel <- link
		}
	}

}

// CountDown will monitor and notify when to close channels
func (cwl *WebCrawler) CountDown() bool {
	count := 0
	for c := range cwl.pendingCountChannel {
		count += c
		// If there are not more pending, then Close
		if count == 0 {
			return true
		}
	}
	return false
}
