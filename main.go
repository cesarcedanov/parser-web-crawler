package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"sync"
	"time"
)

var crawler WebCrawler

func main() {
	crawler = NewWebCrawler()
	fmt.Println("Hello Parser!")
	fmt.Printf("Started at: %s\n", time.Now())
	RunCrawler()
	fmt.Printf("Finished at: %s\n", time.Now())
}

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
			links := crawlLink(crawledLinksChannel, pendingURLChannel, pendingCountChannel)

			crawler.mx.Lock()
			defer crawler.mx.Unlock()
			crawler.crawledLinks = append(crawler.crawledLinks, links...)
		}(w)

	}

	wg.Wait()
	fmt.Printf("%dx Workers found a total of %d valid Links:\n%+v\n\n", crawler.nWorkers, len(crawler.crawledLinks), crawler.crawledLinks)
}

func crawlLink(crawledLinksChannel, pendingLinkChannel chan string, pendingCountChannel chan int) []string {
	links := []string{}
	for link := range pendingLinkChannel {
		inspectURLContent(link, crawledLinksChannel)
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

// initClient will create our Client to do requests
// timeout is 100 secs
func initClient() *http.Client {
	return &http.Client{
		Timeout: 100 * time.Second,
	}
}

// getContentFromURL get the URLs content
func getContentFromURL(url string) (*http.Response, error) {
	response, err := crawler.client.Get(url)
	if err != nil {
		fmt.Printf("Error while getting the Content from URL (%s) due to: %s", url, err)
		return nil, err
	}
	return response, nil
}

// inspectURLContent will connect to the URL
// and inspect the HTML Content to extract link
func inspectURLContent(url string, crawledLinksChannel chan string) {
	response, err := getContentFromURL(url)
	if err != nil {
		fmt.Printf("%s - Error: %s", "inpectURLContent", err)
		return
	}

	defer response.Body.Close()

	z := html.NewTokenizer(response.Body)
	//defer wg.Done()
	for {
		tokenType := z.Next()
		if tokenType == html.ErrorToken {
			return
		}
		token := z.Token()

		if isStartAnchorTag(token, tokenType) {
			link := extractLinkFromTag(token)

			// Append into the queue of link
			// Then send it to the Channel to Crawl them too
			if link != "" {
				go func() {
					crawledLinksChannel <- link
				}()
			}
		}
	}
}

// isStartAnchorTag return true when the token is <a>
func isStartAnchorTag(token html.Token, tokenType html.TokenType) bool {
	//  We need the <a> token instead of the </a>, because those have the href="" attr
	return tokenType == html.StartTagToken && token.DataAtom.String() == "a"
}

// extractLinkFromTag get the href value from the Tag
func extractLinkFromTag(token html.Token) string {
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			if link, isValid := validateLink(crawler.initialURL, attr.Val); isValid {
				return link
			} else {
				// OMIT - Do nothing
			}
		}
	}
	return ""
}

// validateLink check the URL is related to the base URL
func validateLink(base, newURL string) (string, bool) {
	base = strings.TrimSuffix(base, "/")
	switch {
	case strings.HasPrefix(newURL, base):
		return newURL, true
	// I got /career as a newURL (without baseURL) and It should be valid
	case strings.HasPrefix(newURL, "/"):
		return base + newURL, true
	}
	return newURL, false

}
