package client

import (
	"crypto/tls"
	"fmt"
	"github.com/seccijr/quintocrawl/parser"
	"net/http"
	"net/url"
	"github.com/jackdanger/collectlinks"
)

type State struct {
	url    string
	status bool
}

func visitedMonitor() (<-chan map[string]bool, chan <- *State) {
	updates := make(chan *State)
	requests := make(chan map[string]bool)
	urlStatus := make(map[string]bool)
	go func() {
		for {
			select {
			case requests <- urlStatus:
			case state := <-updates:
				urlStatus[state.url] = state.status
			}
		}
	}()

	return requests, updates
}

func Page(url string) {
	queue := make(chan string)
	requests, updates := visitedMonitor()

	go func() { queue <- url }()

	for uri := range queue {
		enqueue(uri, queue, requests, updates)
	}
}

func enqueue(uri string, queue chan string, requests <-chan map[string]bool, updates chan <- *State) {
	fmt.Println("Fetching", uri)
	visited := <-requests
	if (visited[uri]) {
		return
	}
	updates <- &State{url: uri, status: true}
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

	url, err := url.Parse(uri)
	if (err != nil) {
		return
	}

	links := parser.Host(url.Host, resp.Body)

	for _, link := range links {
		absolute := parser.FixUrl(link, uri)
		if uri != "" && !visited[absolute] && absolute != uri {
			go func() {queue <- absolute}()
		}
	}
}
