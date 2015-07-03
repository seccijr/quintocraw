package page

import (
	"io"
	"github.com/PuerkitoBio/goquery"
)

const STATE_SELECT = ".zonas.clearfix a.bold"
const PROV_SELECT = ".zonas.clearfix a[class=\"\"]"
const MAIN_ZONE_SELECT = ".mainZones a"

type PCDoc struct {
	dom goquery.Document
}

func NewDocFromReader(body io.Reader) (*PCDoc, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	return &PCDoc{dom: doc}
}

func has(dom goquery.Document, selector string) (bool, error) {
	chd := dom.Has(selector)
	if chd != nil {
		return true, nil
	}

	return false, nil
}

func getLinks(dom goquery.Document, selector string) ([]string, error) {
	links := make([]string)
	dom.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, hasHref := s.Attr("href")
		if hasHref {
			links = append(href)
		}
	})

	return links
}

func (doc *PCDoc) HasStates() (bool, error) {
	return has(doc.dom, STATE_SELECT)
}

func (doc *PCDoc) HasProvinces() (bool, error) {
	return has(doc.dom, PROV_SELECT)
}

func (doc *PCDoc) GetStateLinks() ([]string, error) {
	return getLinks(doc.dom, STATE_SELECT)
}

func (doc *PCDoc) GetProvinceLinks() ([]string, error) {
	return getLinks(doc.dom, PROV_SELECT)
}
