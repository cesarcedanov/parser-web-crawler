package main

import "strings"

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
