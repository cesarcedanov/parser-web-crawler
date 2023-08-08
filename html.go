package main

import "golang.org/x/net/html"

// isStartAnchorTag return true when the token is <a>
func isStartAnchorToken(token html.Token, tokenType html.TokenType) bool {
	//  We need the <a> token instead of the </a>, because those have the href="" attr
	return tokenType == html.StartTagToken && token.DataAtom.String() == "a"
}

// extractLinkFromTag get the href value from the Tag
func extractLinkFromToken(token html.Token, initialURL string) string {
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			if link, isValid := validateLink(initialURL, attr.Val); isValid {
				return link
			}
		}
	}
	return ""
}
