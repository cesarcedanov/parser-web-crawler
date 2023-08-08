package main

import "testing"

func TestValidateLink(t *testing.T) {
	tests := map[string]struct {
		base     string
		newUrl   string
		expected string
		isValid  bool
	}{
		"ok": {
			base:     "https://parserdigital.com/",
			newUrl:   "https://parserdigital.com/careers",
			expected: "https://parserdigital.com/careers",
			isValid:  true,
		},
		"slash prefix": {
			base:     "https://parserdigital.com/",
			newUrl:   "/careers",
			expected: "https://parserdigital.com/careers",
			isValid:  true,
		},
		"different base": {
			base:     "https://parserdigital.com/",
			newUrl:   "https://cesarcedanov.com/careers",
			expected: "https://cesarcedanov.com/careers",
			isValid:  false,
		},
		"another prefix that contains baseURL": {
			base:     "parserdigital.com/",
			newUrl:   "community.parserdigital.com/",
			expected: "community.parserdigital.com/",
			isValid:  false,
		},
		"empty-base": {
			base:     "",
			newUrl:   "https://cesarcedanov.com/",
			expected: "https://cesarcedanov.com/",
			isValid:  false,
		},
		"empty-newUrl": {
			base:     "https://parserdigital.com/",
			newUrl:   "",
			expected: "",
			isValid:  false,
		},
	}

	for tName, test := range tests {
		t.Run(tName, func(t *testing.T) {
			url, valid := validateLink(test.base, test.newUrl)

			if url != test.expected && valid != test.isValid {
				t.Errorf("Output URL don't match: expected %s, got %s", test.expected, url)
			}
			if valid != test.isValid {
				t.Errorf("Output Valid don't match: expected %v, got %v", test.isValid, valid)
			}

		})
	}

}
