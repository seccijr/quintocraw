package page

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type PCDoc struct {
	dom *goquery.Document
}

func NewDocFromReader(body io.Reader) (*PCDoc, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	return &PCDoc{dom: doc}, nil
}

func Has(dom *goquery.Document, selector string) bool {
	return dom.Find(selector).Is(selector)
}

func GetLinks(dom *goquery.Document, selector string) ([]string, error) {
	var links []string
	dom.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, hasHref := s.Attr("href")
		if hasHref {
			links = append(links, href)
		}
	})

	return links, nil
}
