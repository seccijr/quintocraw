package client

import (
	"net/url"
	"io"
	"errors"
)

type Broker interface {
	// Parses the body of a page request for the current broker.
	// It returns the URL addresses that may follow for the current page.
	// The policies defined to choose between which address may be followed
	// are handled internally in the Broker.
	Parse(httpBody io.Reader) ([]*Broker, error)
	// Gets the page to crawl
	URL() string
}

// Checks if the host of and URL matches the hostname passed as the first
// parameter.
func SameHost(name string, href string) error {
	url, err := url.Parse(href)
	if (err != nil) {
		return err
	}

	if (url.Host != name) {
		return errors.New("Hosts do not match")
	}

	return nil
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

