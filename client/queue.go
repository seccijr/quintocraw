package client

import (
	"crypto/tls"
	"fmt"
	"github.com/jackdanger/collectlinks"
	"net/http"
	"net/url"
)

var visited = make(map[string]bool)

func Page(url string) {
	queue := make(chan string)

	go func() { queue <- url }()

	for uri := range queue {
		enqueue(uri, queue)
	}
}

func enqueue(uri string, queue chan string) {
	fmt.Println("fetching", uri)
	visited[uri] = true
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := fixUrl(link, uri)
		if uri != "" {
			if !visited[absolute] {
				go func() { queue <- absolute }()
			}
		}
	}
}

func fixUrl(href, base string) string {
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
