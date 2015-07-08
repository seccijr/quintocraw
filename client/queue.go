package client

const MAX_BROKERS = 100000

type Queue struct {
	broker 	Broker
	urls	chan string
	proxy   *Proxy
}

func fetchAndParse(url string, proxy *Proxy, broker Broker) ([]string, error) {
	var result []string
	body, err := proxy.Fetch(url)
	if err != nil {
		return result, err
	}
	defer body.Close()

	return broker.Parse(body)
}

func (queue *Queue) Handle(url string) {
	go func() {
		urls, _ := fetchAndParse(url, queue.proxy, queue.broker)
		for _, url := range urls {
			queue.Push(queue.broker.Format(url))
		}
	}()
}

func (queue *Queue) Push(url string) {
	go func() {
		queue.urls <- url
	}()
}

func (queue *Queue) Run() {
	for url := range queue.urls {
		queue.Handle(url)
	}
}

func NewQueue(broker Broker) *Queue {
	urls := make(chan string, MAX_BROKERS)
	proxy := NewProxy()
	return &Queue{broker, urls, proxy}
}
