package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
)

type WebCrawler struct {
	client       *http.Client
	initialURL   string
	nWorkers     int
	crawledLinks []string
	mx           sync.Mutex
}

func NewWebCrawler() WebCrawler {
	var url string
	var nWorkers int
	flag.StringVar(&url, "url", "https://parserdigital.com/", "First URL to Crawl")
	flag.IntVar(&nWorkers, "n", 10, "Number of max workers")

	flag.Parse()
	return WebCrawler{
		client:       initClient(),
		initialURL:   url,
		nWorkers:     nWorkers,
		crawledLinks: []string{},
	}
}

func RunCrawler() {
	pendingCountChannel := make(chan int)
	pendingURLChannel := make(chan string)
	crawledLinksChannel := make(chan string)

	go func() {
		crawledLinksChannel <- crawler.initialURL
	}()

	var wg sync.WaitGroup
	go linkHandler(crawledLinksChannel, pendingURLChannel, pendingCountChannel)
	go func() {
		if countDownAndCloseChannels(pendingCountChannel) {
			close(pendingURLChannel)
			close(crawledLinksChannel)
			close(pendingCountChannel)
		}
	}()

	for w := 1; w <= crawler.nWorkers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			defer crawler.mx.Unlock()

			links := crawlLink(crawledLinksChannel, pendingURLChannel, pendingCountChannel)

			crawler.mx.Lock()
			crawler.crawledLinks = append(crawler.crawledLinks, links...)
		}(w)

	}

	wg.Wait()
	fmt.Printf("%dx Workers found a total of %d valid Links:\n%+v\n\n", crawler.nWorkers, len(crawler.crawledLinks), crawler.crawledLinks)
}

func crawlLink(crawledLinksChannel, pendingLinkChannel chan string, pendingCountChannel chan int) []string {
	links := []string{}
	for link := range pendingLinkChannel {
		inspectURLContent(crawler.initialURL, link, crawledLinksChannel)
		pendingCountChannel <- -1
		links = append(links, link)
	}
	return links
}

func linkHandler(crawledLinksChannel, pendingLinkChannel chan string, pendingCountChannel chan int) {
	alreadyCrawled := make(map[string]bool)

	for link := range crawledLinksChannel {
		if !alreadyCrawled[link] {
			// mark as crawled
			alreadyCrawled[link] = true
			//fmt.Println(link)
			// send it to start crawling it
			pendingCountChannel <- 1
			pendingLinkChannel <- link
		}
	}

}

func countDownAndCloseChannels(pendingCountChannel chan int) bool {
	count := 0
	for c := range pendingCountChannel {
		count += c
		// If there are not more pending, then Close
		if count == 0 {
			return true
		}
	}
	return false
}
