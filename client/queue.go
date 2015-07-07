package client

import (
	"crypto/tls"
	"net/http"
)

const MAX_BROKERS = 64

type Queue struct {
	brokers  chan Broker
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

func enqueue(broker Broker, queue *Queue) {
	uri := broker.URL()
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

	brokers, _ := broker.Parse(resp.Body)
	resp.Body.Close()

	for _, subroker := range brokers {
		newUri := subroker.URL()
		if !visited[newUri] && newUri != uri {
			queue.brokers <- subroker
		}
	}
}

func New() *Queue {
	brokers := make(chan Broker, MAX_BROKERS)
	requests, updates := visitedMonitor()
	q := &Queue{brokers: brokers, requests: requests, updates: updates}

	return q
}

func (queue *Queue) Handle(broker Broker) {
	go func() { queue.brokers <- broker }()
}

func (queue *Queue) Run() {
	for broker := range queue.brokers {
		go func(splitbroker Broker) {
			enqueue(splitbroker, queue)
		}(broker)
	}
}
