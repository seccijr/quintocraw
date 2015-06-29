package parser

import (
	"code.google.com/p/go.net/html"
	"io"
	"net/url"
	"log"
)

func sameHost(name string, href string) bool {
	url, err := url.Parse(href)
	if (err != nil) {
		log.Fatal(err)
	}

	if (url.Host == name) {
		return true
	}
	return false
}

func Host(name string, httpBody io.Reader) []string {
	links := make([]string, 0)
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" && sameHost(name, attr.Val) {
					links = append(links, attr.Val)
				}
			}
		}
	}
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

