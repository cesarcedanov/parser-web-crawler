package main

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"testing"
)

// Read the Documentation to learn more about atom.Atom
// https://pkg.go.dev/golang.org/x/net@v0.14.0/html/atom#pkg-types

func TestIsStartAnchorToken(t *testing.T) {
	tests := map[string]struct {
		token     html.Token
		tokenType html.TokenType
		expected  bool
	}{
		"ok": {
			token: html.Token{
				DataAtom: atom.Atom(0x1), // a
			},
			tokenType: html.StartTagToken,
			expected:  true,
		},
		"error token type": {
			token: html.Token{
				DataAtom: atom.Atom(0x1), // a
			},
			tokenType: html.ErrorToken,
			expected:  false,
		},
		"p not anchor": {
			token: html.Token{
				DataAtom: atom.Atom(0xc01), // p
			},
			tokenType: html.StartTagToken,
			expected:  false,
		},
		"div not anchor": {
			token: html.Token{
				DataAtom: atom.Atom(0x16b03), // div
			},
			tokenType: html.StartTagToken,
			expected:  false,
		},
		"a but ending token": {
			token: html.Token{
				DataAtom: atom.Atom(0x1), // a
			},
			tokenType: html.EndTagToken,
			expected:  false,
		},
	}

	for tName, test := range tests {
		t.Run(tName, func(t *testing.T) {
			actual := isStartAnchorToken(test.token, test.tokenType)

			if actual != test.expected {
				t.Errorf("Output don't match: expected %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestExtractLinkFromToken(t *testing.T) {
	tests := map[string]struct {
		token    html.Token
		url      string
		expected string
	}{
		"ok": {
			token: html.Token{
				Attr: []html.Attribute{
					{
						Key: "href",
						Val: "https://parserdigital.com/careers",
					}, {
						Key: "fake",
						Val: "https://cesarcedanov.com",
					},
				},
			},
			url:      "https://parserdigital.com/",
			expected: "https://parserdigital.com/careers",
		},
		"wrong token key": {
			token: html.Token{
				Attr: []html.Attribute{
					{
						Key: "wrong-key",
						Val: "https://parserdigital.com/",
					},
				},
			},
			url:      "https://parserdigital.com/",
			expected: "",
		},
		"different base url": {
			token: html.Token{
				Attr: []html.Attribute{
					{
						Key: "href",
						Val: "https://cesarcedanov.com/",
					},
				},
			},
			url:      "https://parserdigital.com/",
			expected: "",
		},
	}
	for tName, test := range tests {
		t.Run(tName, func(t *testing.T) {
			actual := extractLinkFromToken(test.token, test.url)

			if actual != test.expected {
				t.Errorf("Output don't match: expected %v, got %v", test.expected, actual)
			}
		})
	}
}
