package interfaces

import (
	"net/url"
	"io"
)

type Parser interface {
	Page(httpBody io.Reader) []string
}


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

