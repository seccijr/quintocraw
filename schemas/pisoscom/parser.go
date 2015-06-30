package pisoscom

import (
	"io"
	"code.google.com/p/go.net/html"
	"github.com/seccijr/quintocrawl/model"
	"net/url"
	"log"
)

type PCParser struct {
	flats model.FlatRepo
}

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

func (parser *PCParser) Page(httpBody io.Reader) []string {
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
