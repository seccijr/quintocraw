package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type Queue struct {
	brokers  chan Broker
	requests chan map[string]bool
	updates  chan string
}

func visitedMonitor() {
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

func enqueue(broker Broker, queue Queue) {
	uri := broker.URL()
	fmt.Println("Fetching", uri)
	visited := <-queue.requests
	if (visited[uri]) {
		return
	}
	queue.updates <- uri
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	hc := http.Client{Transport: transport}
	resp, err := hc.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	brokers, _ := broker.Parse(resp.Body)

	for subroker := range brokers {
		absolute := FixUrl(subroker.URL(), uri)
		if uri != "" && !visited[absolute] && absolute != uri {
			go func() {queue.brokers <- subroker}()
		}
	}
}

func New() *Queue {
	brokers := make(chan Broker)
	requests, updates := visitedMonitor()
	q := &Queue{brokers: brokers, requests: requests, updates: updates}

	return q
}

func (queue *Queue) Handle(broker Broker) {

	go func() { queue.brokers <- broker }()

	for broker := range queue.brokers {
		enqueue(broker, queue)
	}
}
