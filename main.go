package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"time"
)

const parserURL = "https://parserdigital.com/"

func main() {
	fmt.Println("Hello Parser!")

	inspectURLContent(parserURL)
}

// initClient will create our Client to do requests
//
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
func inspectURLContent(url string) {
	response, err := getContentFromURL(url)
	if err != nil {
		fmt.Printf("%s - Error: %s", "inpectURLContent", err)
		return
	}

	defer response.Body.Close()

	z := html.NewTokenizer(response.Body)
	for {
		tokenType := z.Next()
		if tokenType == html.ErrorToken {
			return
		}

		token := z.Token()
		var links []string
		if isStartAnchorTag(token, tokenType) {
			link := extractLinksFromTag(token)
			links = append(links, link)
			// Append into the queue of link
			// Then send it to the Channel to Crawl them too
			fmt.Println(link)
		}

	}
}

// isStartAnchorTag return true when the token is <a>
func isStartAnchorTag(token html.Token, tokenType html.TokenType) bool {
	//  We need the <a> token instead of the </a>, because those have the href="" attr
	return tokenType == html.StartTagToken && token.DataAtom.String() == "a"
}

// extractLinksFromTag get the href value from the Tag
func extractLinksFromTag(token html.Token) string {
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
