package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"time"
)

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

		if isStartAnchorToken(token, tokenType) {
			link := extractLinkFromToken(token)

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
