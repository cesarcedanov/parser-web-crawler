package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"sync"
	"time"
)

const parserURL = "https://parserdigital.com/"

func main() {
	fmt.Println("Hello Parser!")
	fmt.Printf("Started at: %s", time.Now())
	NewCrawler(parserURL)
	fmt.Printf("Finished at: %s", time.Now())
}

func NewCrawler(initialURL string) {
	pendingCountChannel := make(chan int)
	pendingURLChannel := make(chan string)
	crawledLinksChannel := make(chan string)

	go func() {
		crawledLinksChannel <- initialURL
	}()

	var wg sync.WaitGroup
	go linkHandler(crawledLinksChannel, pendingURLChannel, pendingCountChannel)
	go countDownAndCloseChannels(crawledLinksChannel, pendingURLChannel, pendingCountChannel)

	var maxWorkers = 10
	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)
		go crawlLink(&wg, crawledLinksChannel, pendingURLChannel, pendingCountChannel, w)
	}

	wg.Wait()
}

func crawlLink(wg *sync.WaitGroup, crawledLinksChannel, pendingLinkChannel chan string, pendingCountChannel chan int, workerId int) {
	links := []string{}
	for link := range pendingLinkChannel {
		inspectURLContent(link, crawledLinksChannel)
		pendingCountChannel <- -1
		links = append(links, link)
	}
	fmt.Printf("Worker #%d - Found the following Links: %+v\n", workerId, links)
	wg.Done()
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

func countDownAndCloseChannels(crawledLinksChannel chan string, pendingLinkChannel chan string, pendingCountChannel chan int) {
	count := 0
	for c := range pendingCountChannel {
		count += c
		// If there are not more pending, then Close
		if count == 0 {
			close(pendingLinkChannel)
			close(crawledLinksChannel)
			close(pendingCountChannel)
		}
	}
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
	client := initClient()

	response, err := client.Get(url)
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
			if link, isValid := validateLink(parserURL, attr.Val); isValid {
				return link
			} else {
				// append to Found but NOT VALID

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
