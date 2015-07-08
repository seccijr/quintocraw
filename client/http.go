package client

import (
	"net/http"
	"crypto/tls"
	"io"
	"errors"
)

type Proxy struct {
	requests chan map[string]bool
	updates  chan string
}

func visitedMonitor() (chan map[string]bool, chan string) {
	requests := make(chan map[string]bool)
	updates := make(chan string)
	urlStatus := make(map[string]bool)

	go func() {
		for {
			select {
			case requests <- urlStatus:
			case url := <-updates:
				urlStatus[url] = true
			}
		}
	}()

	return requests, updates
}

func NewProxy() *Proxy {
	requests, updates := visitedMonitor()
	return &Proxy{requests, updates}
}

func (proxy *Proxy) Fetch(url string) (io.ReadCloser, error) {
	var reader io.ReadCloser
	visited := <-proxy.requests
	if (visited[url]) {
		return reader, errors.New("Url already visited")
	}

	proxy.updates <- url
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	hc := http.Client{Transport: transport}
	resp, err := hc.Get(url)
	if err != nil {
		return reader, err
	}

	return resp.Body, nil
}