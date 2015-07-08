package client

import (
	"net/url"
	"io"
)

type Broker interface {
	// Parses the body of a page request for the current broker.
	// It returns the URL addresses that may follow for the current page.
	// The policies defined to choose between which address may be followed
	// are handled internally in the Broker.
	Parse(httpBody io.Reader) ([]string, error)
	// Gets the page to crawl
	Format(string) string
}

// Converts a reference from relative to absolute format.
func FixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)

	return uri.String()
}

