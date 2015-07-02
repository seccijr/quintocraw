package pisoscom

import (
	"io"
	"code.google.com/p/go.net/html"
	"github.com/seccijr/quintocrawl/model"
)

type PCBroker struct {
	flats    model.FlatRepo
	url string
}

func (broker *PCBroker) Parse(httpBody io.Reader) []string {
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
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}
}

func (broker *PCBroker) URL() string {
	return broker.url
}
